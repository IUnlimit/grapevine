package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/illtamer/grapevine/internal/stats"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ProxyHandler 核心代理处理器
func ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		suffix := c.Param("suffix")
		path := c.Param("path")

		if isWebSocket(c.Request) {
			handleWebSocket(c, suffix, path)
			return
		}

		handleHTTP(c, suffix, path)
	}
}

func isWebSocket(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

func handleHTTP(c *gin.Context, suffix, path string) {
	ep, err := DefaultResolver.Resolve(suffix)
	if err != nil {
		stats.Global.Record(suffix, true)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	target, _ := url.Parse(fmt.Sprintf("http://%s:%d", ep.Host, ep.Port))

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = path
			req.URL.RawQuery = c.Request.URL.RawQuery
			req.Host = target.Host

			// 透传原始请求头
			for k, vv := range c.Request.Header {
				for _, v := range vv {
					req.Header.Set(k, v)
				}
			}
			req.Header.Set("X-Forwarded-For", c.ClientIP())
			req.Header.Set("X-Real-IP", c.ClientIP())
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
		},
		Transport: &http.Transport{
			DialContext:           (&net.Dialer{Timeout: 3 * time.Second}).DialContext,
			ResponseHeaderTimeout: 30 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error [%s]: %v", suffix, err)
			stats.Global.Record(suffix, true)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error":"backend unavailable"}`))
		},
		ModifyResponse: func(resp *http.Response) error {
			stats.Global.Record(suffix, resp.StatusCode >= 500)
			return nil
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func handleWebSocket(c *gin.Context, suffix, path string) {
	endpoints, err := DefaultResolver.ResolveAll(suffix)
	if err != nil || len(endpoints) == 0 {
		stats.Global.Record(suffix, true)
		c.JSON(http.StatusBadGateway, gin.H{"error": "no available endpoints"})
		return
	}

	// 升级客户端连接
	clientConn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}
	defer clientConn.Close()

	stats.Global.Record(suffix, false)

	// 尝试连接后端
	var backendConn *websocket.Conn
	var connectedIdx int

	for i, ep := range endpoints {
		backendURL := fmt.Sprintf("ws://%s:%d%s", ep.Host, ep.Port, path)
		if c.Request.URL.RawQuery != "" {
			backendURL += "?" + c.Request.URL.RawQuery
		}

		header := http.Header{}
		for k, vv := range c.Request.Header {
			if strings.EqualFold(k, "Upgrade") || strings.EqualFold(k, "Connection") ||
				strings.EqualFold(k, "Sec-Websocket-Key") || strings.EqualFold(k, "Sec-Websocket-Version") ||
				strings.EqualFold(k, "Sec-Websocket-Extensions") || strings.EqualFold(k, "Sec-Websocket-Protocol") {
				continue
			}
			header[k] = vv
		}

		dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
		conn, _, dialErr := dialer.Dial(backendURL, header)
		if dialErr != nil {
			log.Printf("ws dial failed [%s:%d]: %v", ep.Host, ep.Port, dialErr)
			continue
		}
		backendConn = conn
		connectedIdx = i
		break
	}

	if backendConn == nil {
		clientConn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "no backend available"))
		return
	}
	defer backendConn.Close()

	done := make(chan struct{})

	// 后端 -> 客户端
	go func() {
		defer close(done)
		for {
			msgType, msg, err := backendConn.ReadMessage()
			if err != nil {
				// 后端断开，尝试重连到下一个端点
				newConn := reconnectBackend(endpoints, connectedIdx+1, suffix, path, c.Request)
				if newConn != nil {
					backendConn.Close()
					backendConn = newConn
					continue
				}
				return
			}
			if err := clientConn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	}()

	// 客户端 -> 后端
	go func() {
		for {
			msgType, msg, err := clientConn.ReadMessage()
			if err != nil {
				backendConn.Close()
				return
			}
			if err := backendConn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	}()

	<-done
}

func reconnectBackend(_ interface{}, startIdx int, suffix, path string, req *http.Request) *websocket.Conn {
	allEndpoints, err := DefaultResolver.ResolveAll(suffix)
	if err != nil {
		return nil
	}

	for i := startIdx; i < len(allEndpoints); i++ {
		ep := allEndpoints[i]
		backendURL := fmt.Sprintf("ws://%s:%d%s", ep.Host, ep.Port, path)
		if req.URL.RawQuery != "" {
			backendURL += "?" + req.URL.RawQuery
		}

		dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
		conn, _, err := dialer.Dial(backendURL, nil)
		if err != nil {
			continue
		}
		log.Printf("ws reconnected to %s:%d for suffix %s", ep.Host, ep.Port, suffix)
		return conn
	}
	return nil
}

// StreamCopy 双向流拷贝辅助
func StreamCopy(dst io.Writer, src io.Reader, done chan<- struct{}) {
	io.Copy(dst, src)
	done <- struct{}{}
}

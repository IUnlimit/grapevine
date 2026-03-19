package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
)

type HealthChecker struct {
	interval time.Duration
	done     chan struct{}
}

func NewHealthChecker(interval time.Duration) *HealthChecker {
	return &HealthChecker{
		interval: interval,
		done:     make(chan struct{}),
	}
}

func (hc *HealthChecker) Start() {
	go func() {
		ticker := time.NewTicker(hc.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				hc.checkAll()
			case <-hc.done:
				return
			}
		}
	}()
	log.Printf("health checker started (interval: %s)", hc.interval)
}

func (hc *HealthChecker) Stop() {
	close(hc.done)
}

func (hc *HealthChecker) checkAll() {
	var endpoints []model.Endpoint
	database.DB.Find(&endpoints)

	for _, ep := range endpoints {
		alive := checkEndpoint(ep.Host, ep.Port)
		newStatus := "online"
		if !alive {
			newStatus = "offline"
		}

		if ep.Status != newStatus {
			log.Printf("endpoint %s:%d status changed: %s -> %s", ep.Host, ep.Port, ep.Status, newStatus)
			UpdateEndpointStatus(ep.ID, newStatus)
		}
	}
}

func checkEndpoint(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)

	// 先尝试 HTTP 健康检查
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s/", addr))
	if err == nil {
		resp.Body.Close()
		return true
	}

	// 回退到 TCP 探测
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

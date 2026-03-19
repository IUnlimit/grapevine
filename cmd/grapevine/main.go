package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/config"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/handler"
	"github.com/illtamer/grapevine/internal/middleware"
	"github.com/illtamer/grapevine/internal/proxy"
	"github.com/illtamer/grapevine/internal/stats"
)

func main() {
	cfg := config.Load()

	// 初始化数据库
	database.Init(cfg.DBURL, cfg.LogLevel)

	// 初始化 JWT
	middleware.SetJWTSecret(cfg.JWTSecret)

	// 启动统计收集
	stats.Global.Start()
	defer stats.Global.Stop()

	// 启动健康检查
	hc := proxy.NewHealthChecker(10 * time.Second)
	hc.Start()
	defer hc.Stop()

	// 设置 Gin 模式
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 管理 API
	api := r.Group("/api/admin")
	{
		api.GET("/setup/status", handler.SetupStatus)
		api.POST("/setup", handler.Setup)
		api.POST("/login", handler.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/dashboard", handler.Dashboard)

			auth.GET("/services", handler.ListServices)
			auth.GET("/services/:id", handler.GetService)
			auth.POST("/services", handler.CreateService)
			auth.PUT("/services/:id", handler.UpdateService)
			auth.DELETE("/services/:id", handler.DeleteService)

			auth.GET("/endpoints", handler.ListEndpoints)
			auth.GET("/endpoints/:id", handler.GetEndpoint)
			auth.POST("/endpoints", handler.CreateEndpoint)
			auth.PUT("/endpoints/:id", handler.UpdateEndpoint)
			auth.DELETE("/endpoints/:id", handler.DeleteEndpoint)

			auth.GET("/policies/:service_id", handler.GetPolicy)
			auth.PUT("/policies/:service_id", handler.UpsertPolicy)

			auth.GET("/users", handler.ListUsers)
			auth.POST("/users", handler.CreateUser)
			auth.DELETE("/users/:id", handler.DeleteUser)
		}
	}

	// 准备静态文件服务
	webDist := "web/dist"
	var staticFS fs.FS
	var fileServer http.Handler
	hasStatic := false
	if _, err := os.Stat(webDist); err == nil {
		staticFS = os.DirFS(webDist)
		fileServer = http.FileServer(http.FS(staticFS))
		hasStatic = true
	}

	// NoRoute: 静态文件 → 代理转发 → 404
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 1) /api 开头的未匹配路由直接 404
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// 2) /admin/* 管理后台路由 → SPA fallback
		if strings.HasPrefix(path, "/admin") {
			if hasStatic {
				c.Request.URL.Path = "/"
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// 3) 尝试提供静态文件
		if hasStatic {
			trimmed := strings.TrimPrefix(path, "/")
			if trimmed != "" {
				if f, err := fs.Stat(staticFS, trimmed); err == nil && !f.IsDir() {
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}
		}

		// 4) 尝试代理转发: 从路径中提取 /{suffix}/{path...}
		parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
		suffix := parts[0]
		subPath := "/"
		if len(parts) > 1 {
			subPath = "/" + parts[1]
		}

		if suffix == "" {
			// 根路径 "/" → SPA fallback
			if hasStatic {
				c.Request.URL.Path = "/"
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// 设置参数供中间件和 handler 使用
		c.Params = append(c.Params, gin.Param{Key: "suffix", Value: suffix})
		c.Params = append(c.Params, gin.Param{Key: "path", Value: subPath})

		// 限流 → 熔断 → 代理
		middleware.RateLimitHandler(c)
		if c.IsAborted() {
			return
		}
		middleware.CircuitBreakHandler(c)
		if c.IsAborted() {
			return
		}

		// 代理转发
		proxy.ProxyHandler()(c)

		// 记录熔断结果
		middleware.CircuitBreakRecord(c)
	})

	log.Printf("Grapevine gateway starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

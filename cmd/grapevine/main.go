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
		api.POST("/login", handler.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// Dashboard
			auth.GET("/dashboard", handler.Dashboard)

			// Services
			auth.GET("/services", handler.ListServices)
			auth.GET("/services/:id", handler.GetService)
			auth.POST("/services", handler.CreateService)
			auth.PUT("/services/:id", handler.UpdateService)
			auth.DELETE("/services/:id", handler.DeleteService)

			// Endpoints
			auth.GET("/endpoints", handler.ListEndpoints)
			auth.GET("/endpoints/:id", handler.GetEndpoint)
			auth.POST("/endpoints", handler.CreateEndpoint)
			auth.PUT("/endpoints/:id", handler.UpdateEndpoint)
			auth.DELETE("/endpoints/:id", handler.DeleteEndpoint)

			// Policies
			auth.GET("/policies/:service_id", handler.GetPolicy)
			auth.PUT("/policies/:service_id", handler.UpsertPolicy)

			// Users
			auth.GET("/users", handler.ListUsers)
			auth.POST("/users", handler.CreateUser)
			auth.DELETE("/users/:id", handler.DeleteUser)
		}
	}

	// 前端静态文件服务 (生产环境)
	webDist := "web/dist"
	if _, err := os.Stat(webDist); err == nil {
		staticFS := os.DirFS(webDist)
		fileServer := http.FileServer(http.FS(staticFS))

		// 管理后台路由 - 所有非 API、非代理的请求返回 index.html (SPA)
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// API 和代理请求不走静态文件
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}

			// 尝试提供静态文件
			if f, err := fs.Stat(staticFS, strings.TrimPrefix(path, "/")); err == nil && !f.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			// SPA fallback: 返回 index.html
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	}

	// 代理路由 - 带限流和熔断
	proxyGroup := r.Group("/:suffix")
	proxyGroup.Use(middleware.RateLimit(), middleware.CircuitBreak())
	{
		proxyGroup.Any("/*path", proxy.ProxyHandler())
	}

	log.Printf("Grapevine gateway starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

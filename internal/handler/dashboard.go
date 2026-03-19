package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
	"github.com/illtamer/grapevine/internal/stats"
)

type DashboardResponse struct {
	QPS       int64                `json:"qps"`
	Suffixes  []stats.SuffixSummary `json:"suffixes"`
	Endpoints []EndpointStatus     `json:"endpoints"`
}

type EndpointStatus struct {
	ServiceName string `json:"service_name"`
	Suffix      string `json:"suffix"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Priority    int    `json:"priority"`
	Status      string `json:"status"`
}

func Dashboard(c *gin.Context) {
	var services []model.Service
	database.DB.Preload("Endpoints").Find(&services)

	var endpoints []EndpointStatus
	for _, svc := range services {
		for _, ep := range svc.Endpoints {
			endpoints = append(endpoints, EndpointStatus{
				ServiceName: svc.Name,
				Suffix:      svc.Suffix,
				Host:        ep.Host,
				Port:        ep.Port,
				Priority:    ep.Priority,
				Status:      ep.Status,
			})
		}
	}

	resp := DashboardResponse{
		QPS:       stats.Global.GetQPS(),
		Suffixes:  stats.Global.GetSuffixSummaries(),
		Endpoints: endpoints,
	}
	c.JSON(http.StatusOK, resp)
}

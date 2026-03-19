package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
)

func GetPolicy(c *gin.Context) {
	serviceID := c.Param("service_id")
	var cfg model.Config
	if err := database.DB.Where("service_id = ?", serviceID).First(&cfg).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func UpsertPolicy(c *gin.Context) {
	serviceID := c.Param("service_id")
	var input model.Config
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cfg model.Config
	err := database.DB.Where("service_id = ?", serviceID).First(&cfg).Error
	if err != nil {
		// 创建
		input.ServiceID = 0
		database.DB.Raw("SELECT ?::int", serviceID).Scan(&input.ServiceID)
		if err := database.DB.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, input)
		return
	}

	// 更新
	database.DB.Model(&cfg).Updates(map[string]interface{}{
		"rate_limit":              input.RateLimit,
		"burst":                   input.Burst,
		"circuit_break_threshold": input.CircuitBreakThreshold,
	})
	c.JSON(http.StatusOK, cfg)
}

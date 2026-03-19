package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
)

func ListEndpoints(c *gin.Context) {
	serviceID := c.Query("service_id")
	var endpoints []model.Endpoint
	q := database.DB
	if serviceID != "" {
		q = q.Where("service_id = ?", serviceID)
	}
	if err := q.Order("priority ASC").Find(&endpoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, endpoints)
}

func GetEndpoint(c *gin.Context) {
	id := c.Param("id")
	var ep model.Endpoint
	if err := database.DB.First(&ep, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, ep)
}

func CreateEndpoint(c *gin.Context) {
	var ep model.Endpoint
	if err := c.ShouldBindJSON(&ep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&ep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ep)
}

func UpdateEndpoint(c *gin.Context) {
	id := c.Param("id")
	var ep model.Endpoint
	if err := database.DB.First(&ep, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
		return
	}

	var input model.Endpoint
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&ep).Updates(map[string]interface{}{
		"host":     input.Host,
		"port":     input.Port,
		"priority": input.Priority,
		"status":   input.Status,
	})
	c.JSON(http.StatusOK, ep)
}

func DeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Endpoint{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
)

func ListServices(c *gin.Context) {
	var services []model.Service
	if err := database.DB.Preload("Endpoints").Preload("Config").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func GetService(c *gin.Context) {
	id := c.Param("id")
	var svc model.Service
	if err := database.DB.Preload("Endpoints").Preload("Config").First(&svc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, svc)
}

func CreateService(c *gin.Context) {
	var svc model.Service
	if err := c.ShouldBindJSON(&svc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&svc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, svc)
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")
	var svc model.Service
	if err := database.DB.First(&svc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	var input model.Service
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&svc).Updates(map[string]interface{}{
		"suffix":      input.Suffix,
		"name":        input.Name,
		"description": input.Description,
		"is_active":   input.IsActive,
	})
	c.JSON(http.StatusOK, svc)
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Service{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

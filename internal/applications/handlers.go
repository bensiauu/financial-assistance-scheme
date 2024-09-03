package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateApplication(c *gin.Context) {
	var application models.Application
	if err := c.ShouldBindBodyWithJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create application record in DB"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "application created successfully"})
}
func GetAllApplication(c *gin.Context) {
	var applications []models.Application
	if err := db.DB.Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applications)
}
func GetApplicationByID(c *gin.Context) {
	id := c.Param("id")
	var application models.Applicant

	if err := db.DB.First(&application, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "applicantion not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, application)
}

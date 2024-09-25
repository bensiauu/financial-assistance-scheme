package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/internal/utils"
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

	// Check eligibility using the shared utility function
	eligibleSchemes, err := utils.GetEligibleSchemes(application.ApplicantID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isEligible := false
	for _, scheme := range eligibleSchemes {
		if scheme.ID == application.SchemeID {
			isEligible = true
			break
		}
	}

	if !isEligible {
		c.JSON(http.StatusForbidden, gin.H{"error": "Applicant is not eligible for this scheme"})
		return
	}

	if err := db.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application record in DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application created successfully"})
}
func GetAllApplication(c *gin.Context) {
	var applications []models.Application
	if err := db.DB.Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(applications) == 0 {
		c.JSON(http.StatusOK, []models.Application{})
		return
	}

	c.JSON(http.StatusOK, applications)
}
func GetApplicationByID(c *gin.Context) {
	id := c.Param("id")
	var application models.Application

	if err := db.DB.First(&application, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, application)
}

func UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	var application models.Application

	if err := db.DB.First(&application, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	type updateInput struct {
		Status string `json:"status"`
	}
	var input updateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.DB.Model(&application).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application updated successfully"})
}

func DeleteApplication(c *gin.Context) {
	id := c.Param("id")

	var application models.Application
	if err := db.DB.First(&application, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	result := db.DB.Delete(&application)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete application"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}

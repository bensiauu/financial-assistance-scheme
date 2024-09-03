package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/internal/utils"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
)

func CreateScheme(c *gin.Context) {
	var scheme models.Scheme
	if err := c.ShouldBindBodyWithJSON(&scheme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&scheme).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "scheme created successfully"})
}
func GetAllSchemes(c *gin.Context) {
	var schemes []models.Scheme
	if err := db.DB.Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schemes)
}

func GetEligibleSchemes(c *gin.Context) {
	applicantID := c.Query("applicant")
	var applicant models.Applicant
	if err := db.DB.First(&applicant, "id = ?", applicantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
		return
	}

	var schemes []models.Scheme
	if err := db.DB.Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schemes"})
		return
	}

	eligibleSchemes, err := utils.GetEligibleSchemes(applicantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, eligibleSchemes)
}

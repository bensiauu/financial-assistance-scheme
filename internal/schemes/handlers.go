package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/internal/utils"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	if len(schemes) == 0 {
		c.JSON(http.StatusOK, []models.Scheme{})
		return
	}

	var response []models.SchemeResponse

	for _, scheme := range schemes {
		response = append(response, scheme.ToResponse())
	}

	c.JSON(http.StatusOK, response)
}

func GetSchemeByID(c *gin.Context) {
	id := c.Param("id")
	var scheme models.Scheme

	if err := db.DB.First(&scheme, "id = ?", id).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			c.JSON(http.StatusNotFound, gin.H{"error": "scheme not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve administrator"})
		return
	}

	response := scheme.ToResponse()
	c.JSON(http.StatusOK, response)
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

	response := make([]models.SchemeResponse, 0)

	for _, scheme := range eligibleSchemes {
		response = append(response, scheme.ToResponse())
	}

	c.JSON(http.StatusOK, response)
}

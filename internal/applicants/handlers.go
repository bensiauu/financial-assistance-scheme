package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
)

func CreateApplicant(c *gin.Context) {
	var applicant models.Applicant

	if err := c.ShouldBindBodyWithJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&applicant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record in DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "applicant created successfully"})
}

func GetAllApplicants(c *gin.Context) {
	var applicants []models.Applicant
	if err := db.DB.Find(&applicants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applicants)
}
func GetApplicantByID(c *gin.Context) {
	var applicant models.Applicant
	id := c.Param("id")

	if err := db.DB.Find(&applicant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		return
	}

	c.JSON(http.StatusOK, applicant)
}
func UpdateApplicant(c *gin.Context) {
	type Input struct {
		Name             *string                   `json:"name,omitempty"`
		EmploymentStatus *string                   `json:"emplyment_status,omitempty"`
		Sex              *string                   `json:"sex,omitempty"`
		Household        *[]models.HouseholdMember `json:"household,omitempty"`
	}
	var originalApplicant models.Applicant
	var newApplicant Input
	id := c.Param("id")

	if err := c.ShouldBindBodyWithJSON(&newApplicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Find(&originalApplicant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		return
	}

	updates := make(map[string]interface{})
	if newApplicant.Name != nil {
		updates["name"] = *newApplicant.Name
	}
	if newApplicant.EmploymentStatus != nil {
		updates["employment_status"] = *newApplicant.EmploymentStatus
	}
	if newApplicant.Sex != nil {
		updates["sex"] = *newApplicant.Sex
	}
	if err := db.DB.Model(&originalApplicant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if newApplicant.Household != nil {
		// Clear existing household members
		db.DB.Where("applicant_id = ?", originalApplicant.ID).Delete(&models.HouseholdMember{})

		// Add new household members
		for _, member := range *newApplicant.Household {
			member.ApplicantID = originalApplicant.ID
			if err := db.DB.Create(&member).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "applicant updated successfully"})

}
func DeleteApplicant(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Applicant{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "applicant deleted successfully"})
}

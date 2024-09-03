package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateApplicant(c *gin.Context) {
	var input struct {
		Name             string `json:"name,omitempty"`
		EmploymentStatus string `json:"employment_status,omitempty"`
		Sex              string `json:"sex,omitempty"`
		DateOfBirth      string `json:"date_of_birth,omitempty"`
		LastEmployed     string `json:"last_employed,omitempty"`
		Income           int    `json:"income,omitempty"`
		Household        []struct {
			Name             string `json:"name"`
			Relation         string `json:"relation"`
			DateOfBirth      string `json:"date_of_birth"`
			EmploymentStatus string `json:"employment_status"`
		} `json:"household"`
	}

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date of birth"})
		return
	}

	var lastEmployed *time.Time
	if input.LastEmployed != "" {
		t, err := time.Parse("2006-01-02", input.LastEmployed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid last employed date"})
			return
		}
		lastEmployed = &t
	}

	applicant := models.Applicant{
		Name:             input.Name,
		EmploymentStatus: input.EmploymentStatus,
		Sex:              input.Sex,
		DateOfBirth:      dateOfBirth,
		LastEmployed:     lastEmployed,
		Income:           input.Income,
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&applicant).Error; err != nil {
			return err
		}

		// Create Household members
		for _, h := range input.Household {
			hDOB, err := time.Parse("2006-01-02", h.DateOfBirth)
			if err != nil {
				return fmt.Errorf("invalid date of birth for household member: %s", h.Name)
			}

			householdMember := models.HouseholdMember{
				ApplicantID:      applicant.ID,
				Name:             h.Name,
				Relation:         h.Relation,
				DateOfBirth:      hDOB,
				EmploymentStatus: h.EmploymentStatus,
			}

			if err := tx.Create(&householdMember).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create applicant and household members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "applicant created successfully"})
}

func GetAllApplicants(c *gin.Context) {
	var applicants []models.Applicant
	// Use Preload to load Household members
	if err := db.DB.Preload("Household").Find(&applicants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.ApplicantResponse
	for _, applicant := range applicants {
		response = append(response, models.ApplicantResponse{
			ID:               applicant.ID,
			Name:             applicant.Name,
			EmploymentStatus: applicant.EmploymentStatus,
			Sex:              applicant.Sex,
			DateOfBirth:      applicant.DateOfBirth,
			LastEmployed:     applicant.LastEmployed,
			Income:           applicant.Income,
			Household:        applicant.Household,
		})
	}

	c.JSON(http.StatusOK, response)
}
func GetApplicantByID(c *gin.Context) {
	var applicant models.Applicant
	id := c.Param("id")

	if err := db.DB.Preload("Household").First(&applicant, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
			return

		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ApplicantResponse{
		ID:               applicant.ID,
		Name:             applicant.Name,
		EmploymentStatus: applicant.EmploymentStatus,
		Sex:              applicant.Sex,
		DateOfBirth:      applicant.DateOfBirth,
		LastEmployed:     applicant.LastEmployed,
		Household:        applicant.Household,
	}

	c.JSON(http.StatusOK, response)
}
func UpdateApplicant(c *gin.Context) {
	type Input struct {
		Name             *string                   `json:"name,omitempty"`
		DateOfBirth      *string                   `json:"date_of_birth,omitempty"`
		EmploymentStatus *string                   `json:"employment_status,omitempty"`
		Sex              *string                   `json:"sex,omitempty"`
		LastEmployed     *string                   `json:"last_employed,omitempty"`
		Household        *[]models.HouseholdMember `json:"household,omitempty"`
	}
	var originalApplicant models.Applicant
	var newApplicant Input
	id := c.Param("id")

	if err := c.ShouldBindBodyWithJSON(&newApplicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Preload("Household").First(&originalApplicant, "id = ?", id).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if newApplicant.DateOfBirth != nil {
		parsedDate, err := time.Parse("2006-01-02", *newApplicant.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		updates["date_of_birth"] = parsedDate
	}
	if newApplicant.LastEmployed != nil {
		parsedDate, err := time.Parse("2006-01-02", *newApplicant.LastEmployed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		updates["last_employed"] = parsedDate
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

	result := db.DB.Delete(&models.Applicant{}, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "applicant deleted successfully"})
}

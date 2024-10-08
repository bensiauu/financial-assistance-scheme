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
		Name             string `json:"name" binding:"required"`
		EmploymentStatus string `json:"employment_status,omitempty"`
		Sex              string `json:"sex,omitempty"`
		DateOfBirth      string `json:"date_of_birth,omitempty"`
		LastEmployed     string `json:"last_employed,omitempty"`
		Income           int    `json:"income,omitempty"`
		MaritalStatus    string `json:"marital_status,omitempty"`
		DisabilityStatus string `json:"disability_status,omitempty"`
		NumberOfChildren int    `json:"number_of_children,omitempty"`
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
		MaritalStatus:    input.MaritalStatus,
		DisabilityStatus: input.DisabilityStatus,
		NumberOfChildren: input.NumberOfChildren,
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

	// Check if no applicants were found
	if len(applicants) == 0 {
		c.JSON(http.StatusOK, []models.ApplicantResponse{})
		return
	}

	response := make([]models.ApplicantResponse, 0)
	for _, applicant := range applicants {
		response = append(response, models.ApplicantResponse{
			ID:               applicant.ID,
			Name:             applicant.Name,
			EmploymentStatus: applicant.EmploymentStatus,
			Sex:              applicant.Sex,
			DateOfBirth:      applicant.DateOfBirth,
			LastEmployed:     applicant.LastEmployed,
			Income:           applicant.Income,
			MaritalStatus:    applicant.MaritalStatus,
			DisabilityStatus: applicant.DisabilityStatus,
			NumberOfChildren: applicant.NumberOfChildren,
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

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve applicant"})
		return
	}

	response := models.ApplicantResponse{
		ID:               applicant.ID,
		Name:             applicant.Name,
		EmploymentStatus: applicant.EmploymentStatus,
		Sex:              applicant.Sex,
		DateOfBirth:      applicant.DateOfBirth,
		LastEmployed:     applicant.LastEmployed,
		MaritalStatus:    applicant.MaritalStatus,
		DisabilityStatus: applicant.DisabilityStatus,
		NumberOfChildren: applicant.NumberOfChildren,
		Household:        applicant.Household,
	}

	c.JSON(http.StatusOK, response)
}

type updateApplicantInput struct {
	Name             *string `json:"name,omitempty"`
	DateOfBirth      *string `json:"date_of_birth,omitempty"`
	EmploymentStatus *string `json:"employment_status,omitempty"`
	Sex              *string `json:"sex,omitempty"`
	LastEmployed     *string `json:"last_employed,omitempty"`
	Income           *int    `json:"income,omitempty"`
	MaritalStatus    *string `json:"marital_status,omitempty"`
	DisabilityStatus *string `json:"disability_status,omitempty"`
	NumberOfChildren *int    `json:"number_of_children,omitempty"`
	Household        *[]struct {
		Name             string `json:"name"`
		Relation         string `json:"relation"`
		DateOfBirth      string `json:"date_of_birth"`
		EmploymentStatus string `json:"employment_status"`
	} `json:"household,omitempty"`
}

func checkForApplicantUpdates(newApplicant updateApplicantInput) (map[string]interface{}, error) {
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
			return nil, fmt.Errorf("invalid date format for DateOfBirth")
		}
		updates["date_of_birth"] = parsedDate
	}
	if newApplicant.LastEmployed != nil {
		parsedDate, err := time.Parse("2006-01-02", *newApplicant.LastEmployed)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for LastEmployed")
		}
		updates["last_employed"] = parsedDate
	}
	if newApplicant.MaritalStatus != nil {
		updates["marital_status"] = *newApplicant.MaritalStatus
	}
	if newApplicant.DisabilityStatus != nil {
		updates["disability_status"] = *newApplicant.DisabilityStatus
	}
	if newApplicant.Income != nil {
		updates["income"] = newApplicant.Income
	}
	if newApplicant.NumberOfChildren != nil {
		updates["number_of_children"] = newApplicant.NumberOfChildren
	}

	return updates, nil
}

func UpdateApplicant(c *gin.Context) {

	var originalApplicant models.Applicant
	var newApplicant updateApplicantInput
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

	updates, err := checkForApplicantUpdates(newApplicant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
			hDOB, err := time.Parse("2006-01-02", member.DateOfBirth)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date of birth"})
				return
			}

			householdMember := models.HouseholdMember{
				ApplicantID:      originalApplicant.ID,
				Name:             member.Name,
				Relation:         member.Relation,
				DateOfBirth:      hDOB,
				EmploymentStatus: member.EmploymentStatus,
			}
			if err := db.DB.Create(&householdMember).Error; err != nil {
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

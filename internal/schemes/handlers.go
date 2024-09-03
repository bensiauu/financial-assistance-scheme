package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
)

func CreateScheme(c *gin.Context) {
	var scheme models.Scheme
	if err := c.ShouldBindBodyWithJSON(&scheme); err != nil {
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

type Rule struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Criteria struct {
	Rules []Rule `json:"rules"`
}

func evaluateRule(applicant models.Applicant, rule Rule) bool {
	applicantValue := reflect.ValueOf(applicant).FieldByName(rule.Field).Interface()
	switch rule.Operator {
	case "==":
		return applicantValue == rule.Value
	case ">=":
		return applicantValue.(int) >= int(rule.Value.(float64))
	case "<=":
		return applicantValue.(int) <= int(rule.Value.(float64))
	default:
		return false
	}
}

func isApplicantEligible(applicant models.Applicant, criteria Criteria) bool {
	for _, rule := range criteria.Rules {
		if !evaluateRule(applicant, rule) {
			return false
		}
	}
	return true
}

func GetEligibleSchemes(c *gin.Context) {
	applicantID := c.Query("applicant")
	if applicantID == "" {
		// Return an error if the applicant_id is not provided
		c.JSON(http.StatusBadRequest, gin.H{"error": "applicant_id is required"})
		return
	}

	var applicant models.Applicant
	if err := db.DB.First(&applicant, "id = ?", applicantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		return
	}

	var schemes []models.Scheme
	if err := db.DB.Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var eligibleSchemes []models.Scheme
	for _, scheme := range schemes {
		var criteria Criteria
		if err := json.Unmarshal([]byte(scheme.Criteria), &criteria); err != nil {
			continue
		}

		if isApplicantEligible(applicant, criteria) {
			eligibleSchemes = append(eligibleSchemes, scheme)
		}
	}

	c.JSON(http.StatusOK, eligibleSchemes)
}

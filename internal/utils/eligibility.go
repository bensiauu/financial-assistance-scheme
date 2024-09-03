package utils

import (
	"time"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
)

func GetEligibleSchemes(applicantID string) ([]models.Scheme, error) {
	var applicant models.Applicant
	if err := db.DB.First(&applicant, "id = ?", applicantID).Error; err != nil {
		return nil, err
	}

	var schemes []models.Scheme
	if err := db.DB.Find(&schemes).Error; err != nil {
		return nil, err
	}

	var eligibleSchemes []models.Scheme
	for _, scheme := range schemes {
		if isApplicantEligible(applicant, scheme.Criteria) {
			eligibleSchemes = append(eligibleSchemes, scheme)
		}
	}

	return eligibleSchemes, nil
}

func isApplicantEligible(applicant models.Applicant, criteria models.Criteria) bool {
	for _, rule := range criteria.Rules {
		if !evaluateRule(applicant, rule) {
			return false
		}
	}
	return true
}

func evaluateRule(applicant models.Applicant, rule models.Rule) bool {
	switch rule.Field {
	case "income":
		return compareInts(applicant.Income, rule.Operator, int(rule.Value.(float64)))
	case "employment_status":
		return compareStrings(applicant.EmploymentStatus, rule.Operator, rule.Value.(string))
	case "age":
		applicantAge := calculateAge(applicant.DateOfBirth)
		return compareInts(applicantAge, rule.Operator, int(rule.Value.(float64)))
	case "marital_status":
		return compareStrings(applicant.MaritalStatus, rule.Operator, rule.Value.(string))
	case "disability_status":
		return compareStrings(applicant.DisabilityStatus, rule.Operator, rule.Value.(string))
	case "number_of_children":
		return compareInts(applicant.NumberOfChildren, rule.Operator, int(rule.Value.(float64)))
	// Add more fields as needed
	default:
		return false
	}
}

func calculateAge(dob time.Time) int {
	today := time.Now()
	age := today.Year() - dob.Year()
	if today.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

func compareInts(a int, operator string, b int) bool {
	switch operator {
	case "==":
		return a == b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	case ">":
		return a > b
	case "<":
		return a < b
	default:
		return false
	}
}

func compareStrings(a string, operator string, b string) bool {
	switch operator {
	case "==":
		return a == b
	case "!=":
		return a != b
	default:
		return false
	}
}

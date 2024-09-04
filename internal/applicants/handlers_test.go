package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/applicants"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	testDB.AutoMigrate(&models.Applicant{}, &models.HouseholdMember{})

	db.DB = testDB

	t.Cleanup(func() {
		sqlDB, err := db.DB.DB()
		if err != nil {
			t.Logf("Failed to get database connection: %v", err)
		}

		sqlDB.Exec("DROP TABLE IF EXISTS applicants CASCADE")
		sqlDB.Exec("DROP TABLE IF EXISTS household_members CASCADE")
		sqlDB.Close()
	})

	return testDB
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/applicants", handlers.CreateApplicant)
	router.GET("/api/applicants", handlers.GetAllApplicants)
	router.GET("/api/applicants/:id", handlers.GetApplicantByID)
	router.PUT("/api/applicants/:id", handlers.UpdateApplicant)
	router.DELETE("/api/applicants/:id", handlers.DeleteApplicant)
	return router
}

func TestCreateApplicant(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid input",
			inputJSON: `{
		        "name": "John Doe",
		        "employment_status": "employed",
		        "sex": "male",
		        "date_of_birth": "1990-01-01",
		        "income": 50000,
		        "marital_status": "single",
		        "disability_status": "none",
		        "number_of_children": 0,
		        "household": [
		            {"name": "Jane Doe", "relation": "spouse", "date_of_birth": "1992-01-01", "employment_status": "employed"}
		        ]
		    }`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Invalid date format",
			inputJSON: `{
		        "name": "John Doe",
		        "date_of_birth": "01-01-1990"
		    }`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid date of birth",
		},
		{
			name: "Missing name",
			inputJSON: `{
                "employment_status": "employed",
                "date_of_birth": "1990-01-01"
            }`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name: "Database error",
			inputJSON: `{
                "name": "John Doe",
                "date_of_birth": "1990-01-01"
            }`,
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to create applicant and household members",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			if tt.name == "Database error" {
				db.Exec("DROP TABLE applicants CASCADE;")
			}

			req, _ := http.NewRequest("POST", "/api/applicants", strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "applicant created successfully")
			}
		})
	}
}

func TestGetAllApplicants(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func()
		expectedCode  int
		expectedCount int
	}{
		{
			name: "No applicants",
			setupFunc: func() {
				setupTestDB(t)
			},
			expectedCode:  http.StatusOK,
			expectedCount: 0,
		},
		{
			name: "Multiple applicants",
			setupFunc: func() {

				db := setupTestDB(t)
				applicants := []models.Applicant{
					{
						Name:             "John Doe",
						EmploymentStatus: "employed",
						Sex:              "male",
						DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
						Income:           50000,
						MaritalStatus:    "single",
						DisabilityStatus: "none",
						NumberOfChildren: 0,
					},
					{
						Name:             "Jane Doe",
						EmploymentStatus: "unemployed",
						Sex:              "female",
						DateOfBirth:      time.Date(1992, 1, 1, 0, 0, 0, 0, time.UTC),
						Income:           0,
						MaritalStatus:    "married",
						DisabilityStatus: "none",
						NumberOfChildren: 1,
					},
				}
				db.Create(&applicants)
			},
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/applicants", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response []models.ApplicantResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response))
		})
	}
}

func TestGetApplicantByID(t *testing.T) {

	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created applicant
		expectedCode  int
		expectedError string
	}{
		{
			name: "Applicant exists",
			setupFunc: func() string {
				db := setupTestDB(t)
				applicant := models.Applicant{
					Name:             "John Doe",
					EmploymentStatus: "employed",
					Sex:              "male",
					DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(&applicant)
				return applicant.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Applicant not found",
			setupFunc: func() string {
				return uuid.NewString()
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "applicant not found",
		},
		{
			name: "Database error",
			setupFunc: func() string {
				db := setupTestDB(t)
				applicant := models.Applicant{
					Name:             "John Doe",
					EmploymentStatus: "employed",
					Sex:              "male",
					DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(&applicant)
				db.Exec("DROP TABLE applicants CASCADE;")
				return applicant.ID.String()
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to retrieve applicant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicantID := tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/applicants/"+applicantID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				var response models.ApplicantResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, applicantID, response.ID.String())
			}
		})
	}
}

func TestUpdateApplicant(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created applicant
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid update",
			setupFunc: func() string {
				db := setupTestDB(t)
				applicant := models.Applicant{
					Name:             "John Doe",
					EmploymentStatus: "employed",
					Sex:              "male",
					DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(&applicant)
				return applicant.ID.String()
			},
			inputJSON:     `{"name": "John Updated", "employment_status": "unemployed"}`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Applicant not found",
			setupFunc: func() string {
				return uuid.NewString()
			},
			inputJSON:     `{"name": "John Updated"}`,
			expectedCode:  http.StatusNotFound,
			expectedError: "applicant not found",
		},
		{
			name: "Invalid input",
			setupFunc: func() string {
				db := setupTestDB(t)
				applicant := models.Applicant{
					Name:             "John Doe",
					EmploymentStatus: "employed",
					Sex:              "male",
					DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(&applicant)
				return applicant.ID.String()
			},
			inputJSON:     `{"name": 12345}`, // Invalid JSON
			expectedCode:  http.StatusBadRequest,
			expectedError: "json: cannot unmarshal number into Go struct field updateApplicantInput.name of type string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicantID := tt.setupFunc()

			req, _ := http.NewRequest("PUT", "/api/applicants/"+applicantID, strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "applicant updated successfully")
			}
		})
	}
}

func TestDeleteApplicant(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created applicant
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid delete",
			setupFunc: func() string {
				db := setupTestDB(t)
				applicant := models.Applicant{
					Name:             "John Doe",
					EmploymentStatus: "employed",
					Sex:              "male",
					DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(&applicant)
				return applicant.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Applicant not found",
			setupFunc: func() string {
				return uuid.NewString()
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "applicant not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicantID := tt.setupFunc()

			req, _ := http.NewRequest("DELETE", "/api/applicants/"+applicantID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "applicant deleted successfully")
			}
		})
	}
}

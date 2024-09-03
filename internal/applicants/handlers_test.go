package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/applicants"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=testuser password=password123 dbname=test_db port=5432 sslmode=disable TimeZone=UTC"
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	testDB.AutoMigrate(&models.Applicant{}, &models.HouseholdMember{})

	db.DB = testDB

	t.Cleanup(func() {
		sqlDB, err := testDB.DB()
		if err != nil {
			t.Fatalf("Failed to get database connection: %v", err)
		}

		sqlDB.Exec("DROP DATABASE IF EXISTS test_db")
		sqlDB.Close()
	})

	return testDB
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Group("/api").Group("/applicants").
		POST("/", handlers.CreateApplicant).
		GET("/", handlers.GetAllApplicants).
		GET("/:id", handlers.GetApplicantByID).
		PUT("/:id", handlers.UpdateApplicant).
		DELETE("/:id", handlers.DeleteApplicant)
	return router
}

func TestCreateApplicant(t *testing.T) {
	db := setupTestDB(t)
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
			expectedError: "Key: 'input.Name' Error:Field validation for 'Name' failed on the 'required' tag",
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
	db := setupTestDB(t)
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func()
		expectedCode  int
		expectedCount int
	}{
		{
			name:          "No applicants",
			setupFunc:     func() {},
			expectedCode:  http.StatusOK,
			expectedCount: 0,
		},
		{
			name: "Multiple applicants",
			setupFunc: func() {
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
	db := setupTestDB(t)
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
				return "non-existing-id"
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "applicant not found",
		},
		{
			name: "Database error",
			setupFunc: func() string {
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
	db := setupTestDB(t)
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
				return "non-existing-id"
			},
			inputJSON:     `{"name": "John Updated"}`,
			expectedCode:  http.StatusNotFound,
			expectedError: "applicant not found",
		},
		{
			name: "Invalid input",
			setupFunc: func() string {
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
		{
			name: "Database error on update",
			setupFunc: func() string {
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
			inputJSON:     `{"name": "John Updated"}`,
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to hash new password",
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
	db := setupTestDB(t)
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
				return "non-existing-id"
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "applicant not found",
		},
		{
			name: "Database error on delete",
			setupFunc: func() string {
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
			expectedError: "failed to delete applicant",
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

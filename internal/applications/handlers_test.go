package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/applications"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	testDB.AutoMigrate(&models.Application{}, &models.Applicant{}, &models.Scheme{})

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
		POST("/", handlers.CreateApplication).
		GET("/", handlers.GetAllApplication).
		GET("/:id", handlers.GetApplicationByID).
		PUT("/:id", handlers.UpdateApplication).
		DELETE("/:id", handlers.DeleteApplication)
	return router
}

func TestCreateApplication(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() (string, string) // Returns applicantID and schemeID
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid application",
			setupFunc: func() (string, string) {
				applicant := models.Applicant{
					Name: "John Doe", EmploymentStatus: "employed", Sex: "male",
					DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(&applicant)

				scheme := models.Scheme{
					Name: "Low Income Assistance",
					Criteria: models.Criteria{Rules: []models.Rule{
						{Field: "income", Operator: "<=", Value: 20000},
					}},
				}
				db.Create(&scheme)

				return applicant.ID.String(), scheme.ID.String()
			},
			inputJSON:     `{"applicant_id": "<APPLICANT_ID>", "scheme_id": "<SCHEME_ID>"}`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Applicant not eligible",
			setupFunc: func() (string, string) {
				applicant := models.Applicant{
					Name: "John Doe", EmploymentStatus: "employed", Sex: "male",
					DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					Income:      50000,
				}
				db.Create(&applicant)

				scheme := models.Scheme{
					Name: "Low Income Assistance",
					Criteria: models.Criteria{Rules: []models.Rule{
						{Field: "income", Operator: "<=", Value: 20000},
					}},
				}
				db.Create(&scheme)

				return applicant.ID.String(), scheme.ID.String()
			},
			inputJSON:     `{"applicant_id": "<APPLICANT_ID>", "scheme_id": "<SCHEME_ID>"}`,
			expectedCode:  http.StatusForbidden,
			expectedError: "Applicant is not eligible for this scheme",
		},
		{
			name: "Invalid input",
			setupFunc: func() (string, string) {
				return "invalid-applicant-id", "invalid-scheme-id"
			},
			inputJSON:     `{"applicant_id": "invalid-applicant-id", "scheme_id": "invalid-scheme-id"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "json: cannot unmarshal",
		},
		{
			name: "Database error",
			setupFunc: func() (string, string) {
				db.Exec("DROP TABLE applications CASCADE;")
				return "invalid-applicant-id", "invalid-scheme-id"
			},
			inputJSON:     `{"applicant_id": "invalid-applicant-id", "scheme_id": "invalid-scheme-id"}`,
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Failed to create application record in DB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicantID, schemeID := tt.setupFunc()

			inputJSON := strings.Replace(tt.inputJSON, "<APPLICANT_ID>", applicantID, -1)
			inputJSON = strings.Replace(inputJSON, "<SCHEME_ID>", schemeID, -1)

			req, _ := http.NewRequest("POST", "/api/applications", strings.NewReader(inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "Application created successfully")
			}
		})
	}
}

func TestGetAllApplications(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func()
		expectedCode  int
		expectedCount int
	}{
		{
			name:          "No applications",
			setupFunc:     func() {},
			expectedCode:  http.StatusOK,
			expectedCount: 0,
		},
		{
			name: "Multiple applications",
			setupFunc: func() {
				applications := []models.Application{
					{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "pending"},
					{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "approved"},
				}
				db.Create(&applications)
			},
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/applications", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response []models.Application
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response))
		})
	}
}

func TestUpdateApplication(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created application
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid update",
			setupFunc: func() string {
				application := models.Application{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "pending"}
				db.Create(&application)
				return application.ID.String()
			},
			inputJSON:     `{"status": "approved"}`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Application not found",
			setupFunc: func() string {
				return "non-existing-id"
			},
			inputJSON:     `{"status": "approved"}`,
			expectedCode:  http.StatusNotFound,
			expectedError: "application not found",
		},
		{
			name: "Invalid input",
			setupFunc: func() string {
				application := models.Application{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "pending"}
				db.Create(&application)
				return application.ID.String()
			},
			inputJSON:     `{"status": 12345}`, // Invalid JSON
			expectedCode:  http.StatusBadRequest,
			expectedError: "json: cannot unmarshal number into Go struct field",
		},
		{
			name: "Database error on update",
			setupFunc: func() string {
				application := models.Application{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "pending"}
				db.Create(&application)
				db.Exec("DROP TABLE applications CASCADE;")
				return application.ID.String()
			},
			inputJSON:     `{"status": "approved"}`,
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to update application",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicationID := tt.setupFunc()

			req, _ := http.NewRequest("PUT", "/api/applications/"+applicationID, strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "application updated successfully")
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created application
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid delete",
			setupFunc: func() string {
				application := models.Application{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "pending"}
				db.Create(&application)
				return application.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Application not found",
			setupFunc: func() string {
				return "non-existing-id"
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "application not found",
		},
		{
			name: "Database error on delete",
			setupFunc: func() string {
				application := models.Application{ApplicantID: uuid.New(), SchemeID: uuid.New(), Status: "pending"}
				db.Create(&application)
				db.Exec("DROP TABLE applications CASCADE;")
				return application.ID.String()
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to delete application",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicationID := tt.setupFunc()

			req, _ := http.NewRequest("DELETE", "/api/applications/"+applicationID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "Application deleted successfully")
			}
		})
	}
}

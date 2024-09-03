package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/schemes"
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

	testDB.AutoMigrate(&models.Applicant{}, &models.Scheme{})

	db.DB = testDB

	t.Cleanup(func() {
		sqlDB, err := db.DB.DB()
		if err != nil {
			t.Fatalf("Failed to get database connection: %v", err)
		}

		sqlDB.Exec("DROP table IF EXISTS applicants")
		sqlDB.Exec("DROP table IF EXISTS schemes")
		sqlDB.Close()
	})

	return testDB
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/schemes", handlers.CreateScheme)
	router.GET("/api/schemes", handlers.GetAllSchemes)
	router.GET("/api/schemes/eligible", handlers.GetEligibleSchemes)
	return router
}

func TestCreateScheme(t *testing.T) {

	router := setupRouter()

	tests := []struct {
		name          string
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid scheme creation",
			inputJSON: `{
                "name": "Low Income Assistance",
                "criteria": {
                    "rules": [
                        {"field": "income", "operator": "<=", "value": 20000}
                    ]
                },
                "benefits": {
                    "description": "Provides financial assistance to low-income families.",
                    "amount": 1000
                }
            }`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Invalid JSON input",
			inputJSON: `{
                "name": "Low Income Assistance",
                "criteria": "invalid-criteria-format"
            }`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "json: cannot unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			if tt.name == "Database error" {
				db.Exec("DROP TABLE schemes CASCADE;")
			}

			req, _ := http.NewRequest("POST", "/api/schemes", strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "scheme created successfully")
			}
		})
	}
}
func TestGetAllSchemes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func()
		expectedCode  int
		expectedCount int
	}{
		{
			name: "No schemes",
			setupFunc: func() {
				setupTestDB(t)
			},
			expectedCode:  http.StatusOK,
			expectedCount: 0,
		},
		{
			name: "Multiple schemes",
			setupFunc: func() {
				db := setupTestDB(t)
				schemes := []models.Scheme{
					{
						Name: "Low Income Assistance",
						Criteria: models.Criteria{Rules: []models.Rule{
							{Field: "income", Operator: "<=", Value: 20000},
						}},
						Benefits: json.RawMessage(`{"description": "Provides financial assistance to low-income families.", "amount": 1000}`),
					},
					{
						Name: "Housing Assistance",
						Criteria: models.Criteria{Rules: []models.Rule{
							{Field: "housing_status", Operator: "==", Value: "rented"},
						}},
						Benefits: json.RawMessage(`{"description": "Provides financial aid for housing.", "amount": 1500}`),
					},
				}
				db.Create(&schemes)
			},
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/schemes", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response []models.Scheme
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response))
		})
	}
}

func TestGetEligibleSchemes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the applicant ID
		expectedCode  int
		expectedCount int
		expectedError string
	}{
		{
			name: "Eligible schemes found",
			setupFunc: func() string {
				db := setupTestDB(t)
				applicant := models.Applicant{
					Name: "John Doe", EmploymentStatus: "employed", Sex: "male",
					DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					Income:      15000,
				}
				db.Create(&applicant)

				scheme := models.Scheme{
					Name: "Low Income Assistance",
					Criteria: models.Criteria{Rules: []models.Rule{
						{Field: "income", Operator: "<=", Value: 20000},
					}},
					Benefits: json.RawMessage(`{"description": "Provides financial assistance to low-income families.", "amount": 1000}`),
				}
				db.Create(&scheme)

				return applicant.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedCount: 1,
			expectedError: "",
		},
		{
			name: "No eligible schemes found",
			setupFunc: func() string {
				db := setupTestDB(t)
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
					Benefits: json.RawMessage(`{"description": "Provides financial assistance to low-income families.", "amount": 1000}`),
				}
				db.Create(&scheme)

				return applicant.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedCount: 0,
			expectedError: "",
		},
		{
			name: "Applicant not found",
			setupFunc: func() string {
				setupTestDB(t)
				return uuid.NewString()
			},
			expectedCode:  http.StatusNotFound,
			expectedCount: 0,
			expectedError: "Applicant not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicantID := tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/schemes/eligible?applicant="+applicantID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				var response []models.Scheme
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(response))
			}
		})
	}
}

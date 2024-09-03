package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/admin"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=testuser password=password123 dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	testDB.AutoMigrate(&models.Administrator{})

	db.DB = testDB

	t.Cleanup(func() {
		sqlDB, err := db.DB.DB()
		if err != nil {
			t.Logf("Failed to get database connection: %v", err)
		}

		_, err = sqlDB.Exec("DROP TABLE IF EXISTS administrators")
		if err != nil {
			t.Logf("failed to drop db: %v", err)
		}
		sqlDB.Close()
	})

	return testDB
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/administrators", handlers.CreateAdministrator)
	router.GET("/api/administrators", handlers.GetAllAdministrators)
	router.GET("/api/administrators/:id", handlers.GetAdministratorByID)
	router.PUT("/api/administrators/:id", handlers.UpdateAdministrator)
	router.DELETE("/api/administrators/:id", handlers.DeleteAdministrator)
	return router
}

func TestCreateAdministrator(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Valid input",
			inputJSON:     `{"name": "John Doe", "email": "john@example.com", "password": "password123"}`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Missing name",
			inputJSON:     `{"email": "john@example.com", "password": "password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "Key: 'CreateAdminRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name:          "Missing email",
			inputJSON:     `{"name": "John Doe", "password": "password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "Key: 'CreateAdminRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag",
		},
		{
			name:          "Invalid email format",
			inputJSON:     `{"name": "John Doe", "email": "invalid-email", "password": "password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: "Field validation for 'Email'",
		},
		{
			name:          "Duplicate email",
			inputJSON:     `{"name": "John Doe", "email": "john@example.com", "password": "password123"}`,
			expectedCode:  http.StatusConflict,
			expectedError: "email is already in use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			// Prepopulate the database for the duplicate email test case
			if tt.name == "Duplicate email" {
				admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
				db.Create(&admin)
			}

			req, _ := http.NewRequest("POST", "/api/administrators", strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}
		})
	}
}

func TestGetAllAdministrators(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func()
		expectedCode  int
		expectedCount int
	}{
		{
			name:          "No administrators",
			setupFunc:     func() {}, // No setup needed
			expectedCode:  http.StatusOK,
			expectedCount: 0,
		},
		{
			name: "Multiple administrators",
			setupFunc: func() {
				// Add multiple administrators
				admins := []models.Administrator{
					{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"},
					{Name: "Jane Doe", Email: "jane@example.com", PasswordHash: "hashedpassword"},
				}
				db.Create(&admins)
			},
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/administrators", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response []models.AdministratorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response))
		})
	}
}

func TestGetAdministratorByID(t *testing.T) {

	router := setupRouter()

	tests := []struct {
		name          string
		adminID       string
		setupFunc     func() string // Returns the ID of the created admin
		expectedCode  int
		expectedError string
	}{
		{
			name: "Administrator exists",
			setupFunc: func() string {
				db := setupTestDB(t)
				// Add an administrator
				admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
				db.Create(&admin)
				return admin.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Administrator does not exist",
			setupFunc: func() string {
				// No admin setup
				return uuid.New().String()
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "administrator not found",
		},
		{
			name: "Database error",
			setupFunc: func() string {
				db := setupTestDB(t)
				// Simulate a database error by closing the DB connection
				db.Exec("DROP TABLE administrators CASCADE;")
				return "non-existing-id"
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to retrieve administrator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminID := tt.setupFunc()

			req, _ := http.NewRequest("GET", "/api/administrators/"+adminID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				var response models.AdministratorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, adminID, response.ID.String())
			}
		})
	}
}

func TestUpdateAdministrator(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created admin
		inputJSON     string
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid update",
			setupFunc: func() string {
				db := setupTestDB(t)
				// Add an administrator
				admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
				db.Create(&admin)
				return admin.ID.String()
			},
			inputJSON:     `{"name": "Jane Doe", "email": "jane@example.com", "password": "newpassword123"}`,
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Administrator not found",
			setupFunc: func() string {
				// No admin setup
				return uuid.NewString()
			},
			inputJSON:     `{"name": "Jane Doe"}`,
			expectedCode:  http.StatusNotFound,
			expectedError: "administrator not found",
		},
		{
			name: "Invalid input",
			setupFunc: func() string {
				db := setupTestDB(t)
				// Add an administrator
				admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
				db.Create(&admin)
				return admin.ID.String()
			},
			inputJSON:     `{"name": 12345}`, // Invalid JSON
			expectedCode:  http.StatusBadRequest,
			expectedError: "json: cannot unmarshal number into Go struct field Input.name of type string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminID := tt.setupFunc()

			req, _ := http.NewRequest("PUT", "/api/administrators/"+adminID, strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "administrator updated successfully")
			}
		})
	}
}

func TestDeleteAdministrator(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name          string
		setupFunc     func() string // Returns the ID of the created admin
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid delete",
			setupFunc: func() string {
				db := setupTestDB(t)
				// Add an administrator
				admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
				db.Create(&admin)
				return admin.ID.String()
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "Administrator not found",
			setupFunc: func() string {
				// No admin setup
				return uuid.NewString()
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "administrator not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminID := tt.setupFunc()

			req, _ := http.NewRequest("DELETE", "/api/administrators/"+adminID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			} else {
				assert.Contains(t, w.Body.String(), "administrator deleted successfully")
			}
		})
	}
}

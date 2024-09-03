package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/applicants"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDBWithCleanup(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
	db := setupTestDBWithCleanup(t)

	router := setupRouter()

	applicantJSON := `{
		"name": "John Doe",
		"employment_status": "employed",
		"sex": "male",
		"date_of_birth": "1990-01-01",
		"income": 50000,
		"marital_status": "single",
		"disability_status": "none",
		"number_of_children": 0,
		"household": [
			{
				"name": "Jane Doe",
				"relation": "spouse",
				"date_of_birth": "1992-01-01",
				"employment_status": "employed"
			}
		]
	}`

	req, _ := http.NewRequest("POST", "/api/applicants", strings.NewReader(applicantJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "applicant created successfully")

	var applicant models.Applicant
	db.Preload("Household").First(&applicant)
	assert.Equal(t, "John Doe", applicant.Name)
	assert.Equal(t, 1, len(applicant.Household))
	assert.Equal(t, "Jane Doe", applicant.Household[0].Name)
}

func TestGetAllApplicants(t *testing.T) {
	setupTestDBWithCleanup(t)

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/applicants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}

func TestGetApplicantByID(t *testing.T) {
	db := setupTestDBWithCleanup(t)

	router := setupRouter()
	router.GET("/api/applicants/:id", handlers.GetApplicantByID)

	var applicant models.Applicant
	db.First(&applicant)

	req, _ := http.NewRequest("GET", "/api/applicants/"+applicant.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}

func TestUpdateApplicant(t *testing.T) {
	db := setupTestDBWithCleanup(t)

	router := setupRouter()

	var applicant models.Applicant
	db.First(&applicant)

	updateJSON := `{
		"name": "John Updated",
		"date_of_birth": "1991-02-02"
	}`

	req, _ := http.NewRequest("PUT", "/api/applicants/"+applicant.ID.String(), strings.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "applicant updated successfully")

	db.First(&applicant)
	assert.Equal(t, "John Updated", applicant.Name)
}

func TestDeleteApplicant(t *testing.T) {
	db := setupTestDBWithCleanup(t)

	router := setupRouter()

	var applicant models.Applicant
	db.First(&applicant)

	req, _ := http.NewRequest("DELETE", "/api/applicants/"+applicant.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "applicant deleted successfully")

	var count int64
	db.Model(&models.Applicant{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

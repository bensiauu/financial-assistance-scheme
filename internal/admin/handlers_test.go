package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlers "github.com/bensiauu/financial-assistance-scheme/internal/admin"
	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	testDB.AutoMigrate(&models.Administrator{})

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
	router.POST("/api/administrators", handlers.CreateAdministrator)
	router.GET("/api/administrators", handlers.GetAllAdministrators)
	router.GET("/api/administrators/:id", handlers.GetAdministratorByID)
	router.PUT("/api/administrators/:id", handlers.UpdateAdministrator)
	router.DELETE("/api/administrators/:id", handlers.DeleteAdministrator)
	return router
}

func TestCreateAdministrator(t *testing.T) {
	db := setupTestDB(t)

	router := setupRouter()

	adminJSON := `{"name": "John Doe", "email": "john@example.com", "password_hash": "password123"}`
	req, _ := http.NewRequest("POST", "/api/administrators", strings.NewReader(adminJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin created successfully")

	var admin models.Administrator
	db.First(&admin)
	assert.Equal(t, "john@example.com", admin.Email)
}

func TestGetAllAdministrators(t *testing.T) {
	db := setupTestDB(t)

	admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
	db.Create(&admin)

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/administrators", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "john@example.com")
}

func TestGetAdministratorByID(t *testing.T) {
	db := setupTestDB(t)

	admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
	db.Create(&admin)

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/administrators/"+admin.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "john@example.com")
}

func TestUpdateAdministrator(t *testing.T) {
	db := setupTestDB(t)

	// Create a test admin
	admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
	db.Create(&admin)

	router := setupRouter()

	updateJSON := `{"name": "Jane Doe"}`
	req, _ := http.NewRequest("PUT", "/api/administrators/"+admin.ID.String(), strings.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "administrator updated successfully")

	db.First(&admin)
	assert.Equal(t, "Jane Doe", admin.Name)
}

func TestDeleteAdministrator(t *testing.T) {
	db := setupTestDB(t)

	// Create a test admin
	admin := models.Administrator{Name: "John Doe", Email: "john@example.com", PasswordHash: "hashedpassword"}
	db.Create(&admin)

	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/api/administrators/"+admin.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "administrator deleted successfully")

	var count int64
	db.Model(&models.Administrator{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

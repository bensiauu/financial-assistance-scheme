package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateAdministrator(c *gin.Context) {
	var admin models.Administrator

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	admin.PasswordHash = string(hashedPassword)

	if err := db.DB.Create(&admin).Error; err != nil {
		// Check if the error is a PostgreSQL unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "email is already in use"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin created successfully"})
}

func GetAllAdministrators(c *gin.Context) {
	var admins []models.Administrator
	if err := db.DB.Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.AdministratorResponse
	for _, admin := range admins {
		response = append(response, models.AdministratorResponse{
			ID:        admin.ID,
			Name:      admin.Name,
			Email:     admin.Email,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)

}

func GetAdministratorByID(c *gin.Context) {
	id := c.Param("id")
	var admin models.Administrator

	if err := db.DB.First(&admin, "id = ?", id).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			c.JSON(http.StatusNotFound, gin.H{"error": "administrator not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve administrator"})
		return
	}

	response := models.AdministratorResponse{
		ID:        admin.ID,
		Name:      admin.Name,
		Email:     admin.Email,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func UpdateAdministrator(c *gin.Context) {
	type Input struct {
		Name         string `json:"name,omitempty"`
		Email        string `json:"email,omitempty"`
		PasswordHash string `json:"password,omitempty"`
	}
	var prevAdmin models.Administrator
	var newAdmin Input
	id := c.Param("id")

	if err := db.DB.Find(&prevAdmin, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "administrator not found"})
		return
	}

	if err := c.ShouldBindBodyWithJSON(&newAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if newAdmin.Name != "" {
		updates["name"] = newAdmin.Name
	}
	if newAdmin.Email != "" {
		updates["email"] = newAdmin.Email
	}
	if newAdmin.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
			return
		}
		updates["password_hash"] = string(hashedPassword)
	}

	if err := db.DB.Model(&prevAdmin).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "administrator updated successfully"})
}

func DeleteAdministrator(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Delete(&models.Administrator{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete administrator"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "administrator not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "administrator deleted successfully"})
}

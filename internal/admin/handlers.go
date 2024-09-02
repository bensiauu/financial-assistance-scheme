package handlers

import (
	"net/http"

	"github.com/bensiauu/financial-assistance-scheme/models"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdministrator(c *gin.Context) {
	var admin models.Administrator

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to has password"})
		return
	}
	admin.PasswordHash = string(hashedPassword)

	if err := db.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record in DB"})
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

	c.JSON(http.StatusOK, admins)

}

func GetAdministratorByID(c *gin.Context) {
	id := c.Param("id")
	var admin models.Administrator

	if err := db.DB.Find(&admin, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "administrator not found"})
		return
	}

	c.JSON(http.StatusOK, admin)
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

	if newAdmin.Name != "" {
		prevAdmin.Name = newAdmin.Name
	}
	if newAdmin.Email != "" {
		prevAdmin.Email = newAdmin.Email
	}
	if newAdmin.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
			return
		}
		prevAdmin.PasswordHash = string(hashedPassword)
	}

	if err := db.DB.Save(&prevAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "administrator updated successfully"})
}

func DeleteAdministrator(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Administrator{}, "id = ?", id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "administrator not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "administrator deleted successfully"})
}

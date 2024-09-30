package router

import (
	"log"
	"strings"

	admin "github.com/bensiauu/financial-assistance-scheme/internal/admin"
	applicants "github.com/bensiauu/financial-assistance-scheme/internal/applicants"
	applications "github.com/bensiauu/financial-assistance-scheme/internal/applications"
	auth "github.com/bensiauu/financial-assistance-scheme/internal/auth"
	"github.com/bensiauu/financial-assistance-scheme/internal/middleware"
	schemes "github.com/bensiauu/financial-assistance-scheme/internal/schemes"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Handle React Router paths, but exclude API routes

	router.NoRoute(func(c *gin.Context) {
		log.Printf("Request to: %s", c.Request.URL.Path)
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"message": "Not Found"})
		} else {
			c.File("./frontend/build/index.html")
		}
	})

	// Apply CORS middleware globally
	// router.Use(middleware.CORSMiddleware())

	// Login route
	router.POST("/login", auth.Login)

	// Serve static files from the React build directory
	router.Static("/static", "./frontend/build/static")

	// Serve index.html for the root route
	router.StaticFile("/", "./frontend/build/index.html")

	// Handle React Router paths (serve index.html for non-API routes)
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})

	// Protected API routes (after login)
	router.Use(middleware.AuthMiddleware())

	// API routes
	api := router.Group("/api")

	// Administrator routes
	adminGroup := api.Group("/admin")
	{
		adminGroup.POST("/", admin.CreateAdministrator)
		adminGroup.GET("/", admin.GetAllAdministrators)
		adminGroup.GET("/:id", admin.GetAdministratorByID)
		adminGroup.PUT("/:id", admin.UpdateAdministrator)
		adminGroup.DELETE("/:id", admin.DeleteAdministrator)
	}

	// Applicant routes
	applicantGroup := api.Group("/applicants")
	{
		applicantGroup.POST("/", applicants.CreateApplicant)
		applicantGroup.GET("/", applicants.GetAllApplicants)
		applicantGroup.GET("/:id", applicants.GetApplicantByID)
		applicantGroup.PUT("/:id", applicants.UpdateApplicant)
		applicantGroup.DELETE("/:id", applicants.DeleteApplicant)
	}

	// Application routes
	applicationGroup := api.Group("/applications")
	{
		applicationGroup.POST("/", applications.CreateApplication)
		applicationGroup.GET("/", applications.GetAllApplication)
		applicationGroup.GET("/:id", applications.GetApplicationByID)
		applicationGroup.PUT("/:id", applications.UpdateApplication)
		applicationGroup.DELETE("/:id", applications.DeleteApplication)
	}

	// Scheme routes
	schemeGroup := api.Group("/schemes")
	{
		schemeGroup.POST("/", schemes.CreateScheme)
		schemeGroup.GET("/", schemes.GetAllSchemes)
		schemeGroup.GET("/:id", schemes.GetSchemeByID)
		schemeGroup.GET("/eligible", schemes.GetEligibleSchemes)
	}

	return router
}

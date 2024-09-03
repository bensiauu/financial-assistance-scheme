package router

import (
	admin "github.com/bensiauu/financial-assistance-scheme/internal/admin"
	applicants "github.com/bensiauu/financial-assistance-scheme/internal/applicants"
	applications "github.com/bensiauu/financial-assistance-scheme/internal/applications"
	schemes "github.com/bensiauu/financial-assistance-scheme/internal/schemes"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Admin routes
	adminRoutes := router.Group("/admin")
	{
		adminRoutes.POST("/", admin.CreateAdministrator)
		adminRoutes.GET("/", admin.GetAllAdministrators)
		adminRoutes.GET("/:id", admin.GetAdministratorByID)
		adminRoutes.PUT("/:id", admin.UpdateAdministrator)
		adminRoutes.DELETE("/:id", admin.DeleteAdministrator)
	}

	// Applicant routes
	applicantRoutes := router.Group("/applicants")
	{
		applicantRoutes.POST("/", applicants.CreateApplicant)
		applicantRoutes.GET("/", applicants.GetAllApplicants)
		applicantRoutes.GET("/:id", applicants.GetApplicantByID)
		applicantRoutes.PUT("/:id", applicants.UpdateApplicant)
		applicantRoutes.DELETE("/:id", applicants.DeleteApplicant)
	}

	applicationRoutes := router.Group("/applications")
	{
		applicationRoutes.POST("/", applications.CreateApplication)
		applicationRoutes.GET("/", applications.GetAllApplication)
		applicationRoutes.GET("/:id", applications.GetApplicationByID)
	}

	// Scheme routes
	schemeRoutes := router.Group("/schemes")
	{
		schemeRoutes.POST("/", schemes.CreateScheme)
		schemeRoutes.GET("/", schemes.GetAllSchemes)
		schemeRoutes.GET("/eligible/", schemes.GetEligibleSchemes)
	}

	return router
}

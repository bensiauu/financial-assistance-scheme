package router

import (
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

	router.POST("/login", auth.Login)

	router.Use(middleware.AuthMiddleware())
	router.Group("/api").Group("/admin").
		POST("/", admin.CreateAdministrator).
		GET("/", admin.GetAllAdministrators).
		GET("/:id", admin.GetAdministratorByID).
		PUT("/:id", admin.UpdateAdministrator).
		DELETE("/:id", admin.DeleteAdministrator)

	router.Group("/api").Group("/applicants").
		POST("/", applicants.CreateApplicant).
		GET("/", applicants.GetAllApplicants).
		GET("/:id", applicants.GetApplicantByID).
		PUT("/:id", applicants.UpdateApplicant).
		DELETE("/:id", applicants.DeleteApplicant)

	router.Group("/api").Group("/applications").
		POST("/", applications.CreateApplication).
		GET("/", applications.GetAllApplication).
		GET("/:id", applications.GetApplicationByID).
		PUT("/:id", applications.UpdateApplication).
		DELETE("/:id", applications.DeleteApplication)

	router.Group("/api").Group("/schemes").
		POST("/", schemes.CreateScheme).
		GET("/", schemes.GetAllSchemes).
		GET("/eligible/", schemes.GetEligibleSchemes)

	return router
}

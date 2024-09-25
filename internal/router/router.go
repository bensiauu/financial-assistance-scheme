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
	router.Use(middleware.CORSMiddleware())

	router.POST("/login", auth.Login)
	// Serve static files from the React build directory
	router.Static("/static", "./frontend/build/static")

	// Serve the index.html file for the root route
	router.StaticFile("/", "./frontend/build/index.html")

	// For all other routes, serve index.html (to support React Router in client-side routing)
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})

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
		GET("/:id", schemes.GetSchemeByID).
		GET("/eligible/", schemes.GetEligibleSchemes)

	return router
}

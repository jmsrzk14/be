package main

import (
	"fmt"
	"log"
	"os"

	"bem_be/internal/auth"
	"bem_be/internal/auth/campus"
	"bem_be/internal/database"
	"bem_be/internal/handlers"
	"bem_be/internal/middleware"
	"bem_be/internal/services"
	"bem_be/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Set Gin mode
	gin.SetMode(utils.GetEnvWithDefault("GIN_MODE", "debug"))

	// Initialize database connection
	database.Initialize()

	// Initialize auth service (includes both user and student repositories)
	auth.Initialize()
	campusAuthService := services.NewCampusAuthService()

	// Create admin user
	err = auth.CreateAdminUser()
	if err != nil {
		log.Fatalf("Error creating admin user: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	router.Static("/associations", "./uploads/associations")
	router.Static("/clubs", "./uploads/clubs")
	router.Static("/departments", "./uploads/departments")
	router.Static("/bems", "./uploads/bems")
	router.Static("/users", "./uploads/user")

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // []string{"*"}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization", "Content-Type")
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// Register authentication routes
	router.POST("/api/auth/login", handlers.Login)
	router.POST("/api/auth/refresh", handlers.RefreshToken)

	// Login for Student or All Role from External API
	router.POST("/api/auth/campus/login", handlers.CampusLogin)

	// Create handlers
	campusAuthHandler := handlers.NewCampusAuthHandler()

	newsHandler := handlers.NewNewsHandler(database.DB)
	studentHandler := handlers.NewStudentHandler(database.DB, campusAuthService)
	associationHandler := handlers.NewAssociationHandler(database.DB)
	bemHandler := handlers.NewBemHandler(database.DB)
	announcementHandler := handlers.NewAnnouncementHandler(database.DB)
	clubHandler := handlers.NewClubHandler(database.DB)
	galeryHandler := handlers.NewGaleryHandler(database.DB)
	departmentHandler := handlers.NewDepartmentHandler(database.DB)
	organizationHandler := handlers.NewOrganizationHandler(database.DB)
	visimisiHandler := handlers.NewVisiMisiHandler(database.DB)
	requestHandler := handlers.NewRequestHandler(database.DB)

	// Guest Page
	router.GET("/api/association", associationHandler.GetAllAssociationsGuest)
	router.GET("/api/club", clubHandler.GetAllClubsGuest)
	router.GET("/api/department", departmentHandler.GetAllDepartmentsGuest)
	router.GET("/api/bems/manage/:period", bemHandler.GetBEMByPeriod)
	router.GET("/api/visimisibem/:period", visimisiHandler.GetVisiMisiByPeriod)

	// Protected routes
	authRequired := router.Group("/api")
	authRequired.Use(campus.CampusAuthMiddleware())
	{
		// Current user
		authRequired.GET("/auth/me", handlers.GetCurrentUser)

		// Admin routes
		adminRoutes := authRequired.Group("/admin")
		adminRoutes.Use(middleware.RoleMiddleware("Admin"))
		{
			// Campus API token management (admin only)
			adminRoutes.GET("/campus/token", campusAuthHandler.GetToken)
			adminRoutes.POST("/campus/token/refresh", campusAuthHandler.RefreshToken)

			adminRoutes.GET("/organizations/:id", organizationHandler.GetOrganizationByID)

			// Admin access to student data
			adminRoutes.GET("/students", studentHandler.GetAllStudents)
			adminRoutes.GET("/students/:id", studentHandler.GetStudentByID)
			adminRoutes.GET("/students/by-user-id/:user_id", studentHandler.GetStudentByUserID)
			adminRoutes.POST("/students/sync", studentHandler.SyncStudents)
			adminRoutes.PUT("/students/:id/assign", studentHandler.AssignStudent)

			adminRoutes.GET("/news", newsHandler.GetAllNews)
			adminRoutes.GET("/news/:id", newsHandler.GetNewsByID)
			adminRoutes.POST("/news", newsHandler.CreateNews)
			adminRoutes.PUT("/news/:id", newsHandler.UpdateNews)
			adminRoutes.DELETE("/news/:id", newsHandler.DeleteNews)
			adminRoutes.POST("/news/deleted/:id", newsHandler.RestoreNews)

			// Admin access to study program data
			adminRoutes.GET("/clubs", clubHandler.GetAllClubs)
			adminRoutes.GET("/clubs/:id", clubHandler.GetClubByID)
			adminRoutes.POST("/clubs", clubHandler.CreateClub)
			adminRoutes.PUT("/clubs/:id", clubHandler.UpdateClub)
			adminRoutes.DELETE("/clubs/:id", clubHandler.DeleteClub)

			// Admin access to clubassociation data
			adminRoutes.GET("/association", associationHandler.GetAllAssociations)
			adminRoutes.GET("/associations/:id", associationHandler.GetAssociationByID)
			adminRoutes.POST("/associations", associationHandler.CreateAssociation)
			adminRoutes.PUT("/associations/:id", associationHandler.UpdateAssociation)
			adminRoutes.DELETE("/associations/:id", associationHandler.DeleteAssociation)

			adminRoutes.GET("/bem", bemHandler.GetAllBems)
			adminRoutes.GET("/bems/:id", bemHandler.GetBemByID)
			adminRoutes.POST("/bems", bemHandler.CreateBem)
			adminRoutes.PUT("/bems/:id", bemHandler.UpdateBem)
			adminRoutes.DELETE("/bems/:id", bemHandler.DeleteBem)

			adminRoutes.GET("/announcement", announcementHandler.GetAllAnnouncements)
			adminRoutes.GET("/announcements/:id", announcementHandler.GetAnnouncementByID)
			adminRoutes.POST("/announcements", announcementHandler.CreateAnnouncement)
			adminRoutes.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)
			adminRoutes.DELETE("/announcements/:id", announcementHandler.DeleteAnnouncement)

			adminRoutes.GET("/galery", galeryHandler.GetAllGalerys)
			adminRoutes.GET("/galery/:id", galeryHandler.GetGaleryByID)
			adminRoutes.POST("/galery", galeryHandler.CreateGalery)
			adminRoutes.PUT("/galery/:id", galeryHandler.UpdateGalery)
			adminRoutes.DELETE("/galery/:id", galeryHandler.DeleteGalery)

			adminRoutes.GET("/department", departmentHandler.GetAllDepartments)
			adminRoutes.GET("/department/:id", departmentHandler.GetDepartmentByID)
			adminRoutes.POST("/department", departmentHandler.CreateDepartment)
			adminRoutes.PUT("/department/:id", departmentHandler.UpdateDepartment)
			adminRoutes.DELETE("/department/:id", departmentHandler.DeleteDepartment)

			adminRoutes.GET("/request", requestHandler.GetAllRequests)
			adminRoutes.GET("/request/:id", requestHandler.GetRequestByID)
			adminRoutes.POST("/request", requestHandler.CreateRequest)
			adminRoutes.PUT("/request/:id", requestHandler.UpdateRequest)
			adminRoutes.DELETE("/request/:id", requestHandler.DeleteRequest)
		}

		// Employee routes (replacing assistant routes)
		studentRoutes := authRequired.Group("/student")
		studentRoutes.Use(middleware.RoleMiddleware("Mahasiswa"))
		{
			studentRoutes.GET("/visimisibem/:id", visimisiHandler.GetVisiMisiById)
			studentRoutes.PUT("/visimisibem/:id", visimisiHandler.UpdateVisiMisiBem)
			studentRoutes.PUT("/visimisiperiod/:id", visimisiHandler.UpdateVisiMisiPeriod)
			studentRoutes.POST("/announcements", announcementHandler.CreateAnnouncement)

			studentRoutes.GET("/clubs", clubHandler.GetAllClubs)
			studentRoutes.GET("/clubs/:id", clubHandler.GetClubByID)

			studentRoutes.GET("/departments", departmentHandler.GetAllDepartments)
			studentRoutes.GET("/departments/:id", departmentHandler.GetDepartmentByID)

			studentRoutes.GET("/associations", associationHandler.GetAllAssociations)
			studentRoutes.GET("/associations/:id", associationHandler.GetAssociationByID)
			studentRoutes.GET("/profile", handlers.GetCurrentUser)
			studentRoutes.PUT("/profile", handlers.EditProfile)
		}

		// Assistant routes
		assistantRoutes := authRequired.Group("/assistant")
		assistantRoutes.Use(middleware.RoleMiddleware("Asisten Dosen", "asisten dosen"))
		{
		}
	}

	// Start the server
	port := utils.GetEnvWithDefault("SERVER_PORT", "9090")

	// Add public endpoints
	router.GET("/api/students/by-user-id/:user_id", studentHandler.GetStudentByUserID)

	// setelah semua route didefinisikan
	for _, ri := range router.Routes() {
		fmt.Println(ri.Method, ri.Path)
	}

	log.Printf("Server running on port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}

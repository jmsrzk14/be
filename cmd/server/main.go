package main

import (
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
)

func main() {
	// Validate required environment variables
	requiredEnvVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
	}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			log.Fatalf("Variabel lingkungan %s tidak diatur", env)
		}
	}

	// Set Gin mode
	gin.SetMode(utils.GetEnvWithDefault("GIN_MODE", "debug"))

	// Initialize database connection
	database.Initialize()

	// Initialize auth service (includes both user and student repositories)
	auth.Initialize()
	campusAuthService := services.NewCampusAuthService()

	// Create admin user
	err := auth.CreateAdminUser()
	if err != nil {
		log.Fatalf("Gagal membuat pengguna admin: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	router.Static("/associations", "./uploads/associations")
	router.Static("/clubs", "./uploads/clubs")
	router.Static("/departments", "./uploads/departments")
	router.Static("/bems", "./uploads/bems")
	router.Static("/users", "./uploads/user")
	router.Static("/requests", "./uploads/requests")
	router.Static("/news", "./Uploads/news")
	router.Static("/request", "./Uploads/requests")
	router.Static("/barang", "./Uploads/request/barang")

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
	itemHandler := handlers.NewItemHandler(database.DB)
	aspirationHandler := *handlers.NewAspirationHandler(database.DB)
	eventHandler := handlers.NewEventHandler(database.DB)

	// Guest Page
	router.GET("/api/association", associationHandler.GetAllAssociationsGuest)
	router.GET("/api/club", clubHandler.GetAllClubsGuest)
	router.GET("/api/department", departmentHandler.GetAllDepartmentsGuest)
	router.GET("/api/bems/manage/:period", bemHandler.GetBEMByPeriod)
	router.GET("/api/visimisibem/:period", visimisiHandler.GetVisiMisiByPeriod)
	router.GET("/api/news", newsHandler.GetAllNews)
	router.GET("/api/announcements", announcementHandler.GetAllAnnouncement)
	router.GET("/api/news/:id", newsHandler.GetNewsByID)
	router.GET("/api/item_sarpras", itemHandler.GetAllItemsSarpras)
	router.GET("/api/item_depol", itemHandler.GetAllItemsDepol)
	router.GET("/api/events", eventHandler.GetEventsCurrentMonth)

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

			// Admin access to association data
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

			adminRoutes.GET("/announcement", announcementHandler.GetAllAnnouncement)
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

			adminRoutes.GET("/request_sarpras", requestHandler.GetAllRequestsSarpras)
			adminRoutes.GET("/request_sarpras/:id", requestHandler.GetRequestByIDSarpras)

			adminRoutes.GET("/request_depol", requestHandler.GetAllRequestsDepol)
			adminRoutes.GET("/request_depol/:id", requestHandler.GetRequestByIDDepol)

			adminRoutes.GET("/item", itemHandler.GetAllItemsSarpras)
			adminRoutes.GET("/item/:id", itemHandler.GetItemSarparsByID)

			adminRoutes.GET("/aspirations", aspirationHandler.GetAllAspirations)
		}

		// Student routes
		studentRoutes := authRequired.Group("/student")
		studentRoutes.Use(middleware.RoleMiddleware("Mahasiswa"))
		{
			studentRoutes.GET("/visimisibem/:id", visimisiHandler.GetVisiMisiById)
			studentRoutes.PUT("/visimisibem/:id", visimisiHandler.UpdateVisiMisiBem)
			studentRoutes.PUT("/visimisiperiod/:id", visimisiHandler.UpdateVisiMisiPeriod)
			studentRoutes.POST("/announcements", announcementHandler.CreateAnnouncement)
			studentRoutes.GET("/announcement", announcementHandler.GetAllAnnouncement)
			studentRoutes.GET("/announcements/:id", announcementHandler.GetAnnouncementByID)
			studentRoutes.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)

			studentRoutes.GET("/news", newsHandler.GetAllNews)
			studentRoutes.GET("/news/:id", newsHandler.GetNewsByID)
			studentRoutes.POST("/news", newsHandler.CreateNews)
			studentRoutes.PUT("/news/:id", newsHandler.UpdateNews)
			studentRoutes.DELETE("/news/:id", newsHandler.DeleteNews)


			studentRoutes.GET("/clubs", clubHandler.GetAllClubs)
			studentRoutes.GET("/clubs/:id", clubHandler.GetClubByID)

			studentRoutes.GET("/departments", departmentHandler.GetAllDepartments)
			studentRoutes.GET("/departments/:id", departmentHandler.GetDepartmentByID)

			studentRoutes.GET("/associations", associationHandler.GetAllAssociations)
			studentRoutes.GET("/associations/:id", associationHandler.GetAssociationByID)
			studentRoutes.GET("/profile", handlers.GetCurrentUser)
			studentRoutes.PUT("/profile", handlers.EditProfile)

			studentRoutes.POST("/requests_sarpras", requestHandler.CreateRequestSarpras)
			studentRoutes.GET("/request_sarpras", requestHandler.GetAllRequestsSarpras)
			studentRoutes.GET("/request_sarpras/:id", requestHandler.GetRequestByIDSarpras)
			studentRoutes.GET("/request_sarpras/user/:id", requestHandler.GetRequestsByUserIDSapras)
			studentRoutes.PUT("/request_sarpras/:id", requestHandler.UpdateRequestSarpras)
			studentRoutes.PUT("/request_sarpras/image_barang/:id", requestHandler.UploadImageBarangSarpras)
			studentRoutes.PUT("/request_sarpras/status/:id", requestHandler.UpdateRequestSarprasStatus)
			studentRoutes.PUT("/request_sarpras/return/:id", requestHandler.ReturnBarangSarpras)
			studentRoutes.PUT("/request_sarpras/done/:id", requestHandler.EndRequestBarangSarpras)
			studentRoutes.DELETE("/request_sarpras/:id", requestHandler.DeleteRequestSarpras)

			studentRoutes.POST("/requests_depol", requestHandler.CreateRequestDepol)
			studentRoutes.GET("/request_depol", requestHandler.GetAllRequestsDepol)
			studentRoutes.GET("/request_depol/:id", requestHandler.GetRequestByIDDepol)
			studentRoutes.GET("/request_depol/user/:id", requestHandler.GetRequestsByUserIDDepol)
			studentRoutes.PUT("/request_depol/:id", requestHandler.UpdateRequestDepol)
			studentRoutes.PUT("/request_depol/image_barang/:id", requestHandler.UploadImageBarangDepol)
			studentRoutes.PUT("/request_depol/status/:id", requestHandler.UpdateRequestDepolStatus)
			studentRoutes.PUT("/request_depol/return/:id", requestHandler.ReturnBarangDepol)
			studentRoutes.PUT("/request_depol/done/:id", requestHandler.EndRequestBarangDepol)
			studentRoutes.DELETE("/request_depol/:id", requestHandler.DeleteRequestDepol)

			studentRoutes.GET("/item_sarpras", itemHandler.GetAllItemsSarpras)
			studentRoutes.GET("/item_sarpras/:id", itemHandler.GetItemSarparsByID)
			studentRoutes.POST("/item_sarpras", itemHandler.CreateItemSarpras)
			studentRoutes.PUT("/item_sarpras/:id", itemHandler.UpdateItemSarpras)
			studentRoutes.DELETE("/item_sarpras/:id", itemHandler.DeleteItemSarpras)

			studentRoutes.GET("/item_depol", itemHandler.GetAllItemsDepol)
			studentRoutes.GET("/item_depol/:id", itemHandler.GetItemDepolByID)
			studentRoutes.POST("/item_depol", itemHandler.CreateItemDepol)
			studentRoutes.PUT("/item_depol/:id", itemHandler.UpdateItemDepol)
			studentRoutes.DELETE("/item_depol/:id", itemHandler.DeleteItemDepol)

			studentRoutes.POST("/aspirations", aspirationHandler.CreateAspiration)

			studentRoutes.POST("/events", eventHandler.CreateEvent)
			studentRoutes.PUT("/events/:id", eventHandler.UpdateEvent)
			studentRoutes.GET("/events/current-month", eventHandler.GetEventsCurrentMonth)
			studentRoutes.DELETE("/events/:id", eventHandler.DeleteEvent)
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

	// // Log all registered routes
	// for _, ri := range router.Routes() {
	// 	log.Printf("Route: %s %s", ri.Method, ri.Path)
	// }

	log.Printf("Server berjalan di port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
		os.Exit(1)
	}
}

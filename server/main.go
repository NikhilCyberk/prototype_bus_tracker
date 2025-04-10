// main.go
package main

import (
	"bustracking/config"
	"bustracking/handlers"
	"bustracking/middleware"
	"bustracking/models"
	"bustracking/seeder"
	"bustracking/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Bus{},
		&models.Route{},
		&models.BusStop{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create initial seed data if needed
	seedData(db)

	// Initialize services
	authService := services.NewAuthService(db)
	// mapsService := services.NewMapsService(cfg.MapsAPIKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	wsHandler := handlers.NewWebSocketHandler()
	busHandler := handlers.NewBusHandler(db, wsHandler.Broadcast)
	routeHandler := handlers.NewRouteHandler(db)
	stopHandler := handlers.NewStopHandler(db)

	// Start WebSocket broadcaster
	go wsHandler.BroadcastMessages()

	// Set up router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Authentication routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Authenticated routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", authHandler.GetProfile)

		// Bus routes
		buses := api.Group("/buses")
		{
			buses.GET("/", busHandler.GetAllBuses)
			buses.GET("/:id", busHandler.GetBusById)
			buses.PUT("/:id/location", busHandler.UpdateBusLocation)
			buses.PUT("/:id/status", busHandler.UpdateBusStatus)
		}
		// Admin-only routes
		adminBusRoutes := buses.Group("/")
		adminBusRoutes.Use(middleware.RoleMiddleware("admin"))
		{
			adminBusRoutes.POST("/", busHandler.CreateBus)
			adminBusRoutes.PUT("/:id", busHandler.UpdateBus)
			adminBusRoutes.DELETE("/:id", busHandler.DeleteBus)
			adminBusRoutes.POST("/assign", busHandler.AssignBusToDriver)
		}

		// Route information
		routes := api.Group("/routes")
		{
			routes.GET("/", routeHandler.GetAllRoutes)
			routes.GET("/:id", routeHandler.GetRouteById)
			routes.GET("/:id/stops", routeHandler.GetRouteStops)
			routes.POST("/", middleware.RoleMiddleware("admin"), routeHandler.CreateRoute)
			routes.PUT("/:id", middleware.RoleMiddleware("admin"), routeHandler.UpdateRoute)
		}

		// Bus stops
		stops := api.Group("/stops")
		{
			stops.GET("/", stopHandler.GetAllStops)
			stops.GET("/:id", stopHandler.GetStopById)
			stops.GET("/:id/arrivals", stopHandler.GetStopArrivals)
			stops.POST("/", middleware.RoleMiddleware("admin"), stopHandler.CreateStop)
			stops.PUT("/:id", middleware.RoleMiddleware("admin"), stopHandler.UpdateStop)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.POST("/assign-bus", busHandler.AssignBusToDriver)
		}
	}

	// WebSocket endpoint
	router.GET("/ws", wsHandler.HandleWebSocket)

	// Start server
	log.Println("Server started on " + cfg.GetServerAddress())
	log.Fatal(router.Run(cfg.GetServerAddress()))
}

func seedData(db *gorm.DB) {
	// Create a new seeder
	dataSeed := seeder.NewSeeder(db)

	// Seed the data
	if err := dataSeed.SeedIndianData(); err != nil {
		log.Println("Error seeding data:", err)
	}
}

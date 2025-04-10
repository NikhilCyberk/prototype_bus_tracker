// // main.go
// package main

// import (
// 	"bustracking/seeder"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// // Global variables
// var (
// 	db        *gorm.DB
// 	wsClients = make(map[*websocket.Conn]bool)
// 	broadcast = make(chan Message)
// 	upgrader  = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true // Allow all connections for now
// 		},
// 	}
// )

// // Database models
// type User struct {
// 	gorm.Model
// 	Username string `gorm:"unique" json:"username"`
// 	Password string `json:"-"`
// 	Role     string `json:"role"` // "driver" or "passenger"
// }

// type Bus struct {
// 	gorm.Model
// 	BusNumber   string    `json:"busNumber"`
// 	DriverID    uint      `json:"driverId"`
// 	Driver      User      `json:"driver"`
// 	RouteID     uint      `json:"routeId"`
// 	Route       Route     `json:"route"`
// 	Status      string    `json:"status"` // "on-route", "delayed", "off-duty"
// 	Latitude    float64   `json:"latitude"`
// 	Longitude   float64   `json:"longitude"`
// 	LastUpdated time.Time `json:"lastUpdated"`
// }

// type Route struct {
// 	gorm.Model
// 	RouteNumber string     `json:"routeNumber"`
// 	Name        string     `json:"name"`
// 	Color       string     `json:"color"`
// 	Stops       []BusStop  `json:"stops"`
// 	Path        []GeoPoint `gorm:"serializer:json" json:"path"`
// }

// type BusStop struct {
// 	gorm.Model
// 	Name      string  `json:"name"`
// 	RouteID   uint    `json:"routeId"`
// 	Order     int     `json:"order"`
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

// type GeoPoint struct {
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

// type Message struct {
// 	Type    string      `json:"type"`
// 	Content interface{} `json:"content"`
// }

// type BusLocationUpdate struct {
// 	BusID     uint    `json:"busId"`
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// 	Status    string  `json:"status"`
// 	Timestamp int64   `json:"timestamp"`
// }

// func main() {
// 	// Initialize database
// 	initDB()

// 	// Create initial seed data if needed
// 	seedData()

// 	// Set up router
// 	router := gin.Default()

// 	// Configure CORS
// 	router.Use(cors.New(cors.Config{
// 		// AllowOrigins:     []string{"*"},
// 		AllowOrigins:     []string{"http://localhost:5173"},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	// Set up API routes
// 	setupRoutes(router)

// 	// Start WebSocket handler
// 	go handleMessages()

// 	// Start server
// 	log.Println("Server started on :8080")
// 	log.Fatal(router.Run(":8080"))
// }

// func initDB() {
// 	var err error
// 	dsn := "host=localhost user=postgres password=1234 dbname=bustrackerdb port=5432 sslmode=disable"
// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Auto migrate schemas
// 	db.AutoMigrate(&User{}, &Bus{}, &Route{}, &BusStop{})
// 	log.Println("Database initialized successfully")
// }

// // func seedData() {
// // 	// Check if we need to seed data
// // 	var count int64
// // 	db.Model(&Route{}).Count(&count)
// // 	if count > 0 {
// // 		return // Data already exists
// // 	}

// // 	// Seed some routes and stops
// // 	routes := []Route{
// // 		{
// // 			RouteNumber: "101",
// // 			Name:        "Downtown Express",
// // 			Color:       "#FF0000",
// // 			Path: []GeoPoint{
// // 				{Latitude: 37.7749, Longitude: -122.4194},
// // 				{Latitude: 37.7750, Longitude: -122.4180},
// // 				{Latitude: 37.7755, Longitude: -122.4130},
// // 			},
// // 		},
// // 		{
// // 			RouteNumber: "202",
// // 			Name:        "Uptown Local",
// // 			Color:       "#00FF00",
// // 			Path: []GeoPoint{
// // 				{Latitude: 37.7849, Longitude: -122.4294},
// // 				{Latitude: 37.7850, Longitude: -122.4280},
// // 				{Latitude: 37.7855, Longitude: -122.4230},
// // 			},
// // 		},
// // 	}

// // 	for i := range routes {
// // 		db.Create(&routes[i])

// // 		// Create stops for this route
// // 		for j := 0; j < 5; j++ {
// // 			offset := float64(j) * 0.002
// // 			stop := BusStop{
// // 				Name:      "Stop " + routes[i].RouteNumber + "-" + string(rune(65+j)), // A, B, C, etc.
// // 				RouteID:   routes[i].ID,
// // 				Order:     j + 1,
// // 				Latitude:  routes[i].Path[0].Latitude + offset,
// // 				Longitude: routes[i].Path[0].Longitude + offset,
// // 			}
// // 			db.Create(&stop)
// // 		}
// // 	}

// // 	for i := range routes {
// // 		db.Create(&routes[i])

// // 		// Create stops for this route
// // 		for j := 0; j < 5; j++ {
// // 			offset := float64(j) * 0.002
// // 			stop := BusStop{
// // 				Name:      "Stop " + routes[i].RouteNumber + "-" + string(rune(65+j)), // A, B, C, etc.
// // 				RouteID:   routes[i].Model.ID,                                         // Access the ID field using the Model struct
// // 				Order:     j + 1,
// // 				Latitude:  routes[i].Path[0].Latitude + offset,
// // 				Longitude: routes[i].Path[0].Longitude + offset,
// // 			}
// // 			db.Create(&stop)
// // 		}
// // 	}

// // 	log.Println("Seed data created successfully")
// // }

// func seedData() {
// 	// Create a new seeder
// 	dataSeed := seeder.NewSeeder(db)

// 	// Seed the data
// 	if err := dataSeed.SeedIndianData(); err != nil {
// 		log.Println("Error seeding data:", err)
// 	}
// }

// // Add this handler function
// func assignBusToDriver(c *gin.Context) {
// 	var assignment struct {
// 		DriverID uint `json:"driverId"`
// 		BusID    uint `json:"busId"`
// 	}

// 	if err := c.ShouldBindJSON(&assignment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if driver exists
// 	var driver User
// 	if err := db.First(&driver, assignment.DriverID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
// 		return
// 	}

// 	// Check if driver role is correct
// 	if driver.Role != "driver" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a driver"})
// 		return
// 	}

// 	// Check if bus exists
// 	var bus Bus
// 	if err := db.First(&bus, assignment.BusID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
// 		return
// 	}

// 	// Update bus with driver
// 	bus.DriverID = assignment.DriverID
// 	if err := db.Save(&bus).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign bus"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Bus assigned successfully", "bus": bus})
// }

// func setupRoutes(router *gin.Engine) {
// 	// Authentication routes
// 	auth := router.Group("/api/auth")
// 	{
// 		auth.POST("/register", registerUser)
// 		auth.POST("/login", loginUser)
// 	}

// 	admin := router.Group("/api/admin")
// 	{
// 		// admin.Use(authMiddleware()) // Add authentication middleware in production
// 		admin.POST("/assign-bus", assignBusToDriver)
// 	}

// 	// Bus routes
// 	buses := router.Group("/api/buses")
// 	{
// 		buses.GET("/", getAllBuses)
// 		buses.GET("/:id", getBusById)
// 		buses.PUT("/:id/location", updateBusLocation)
// 		buses.PUT("/:id/status", updateBusStatus)
// 	}

// 	// Route information
// 	routes := router.Group("/api/routes")
// 	{
// 		routes.GET("/", getAllRoutes)
// 		routes.GET("/:id", getRouteById)
// 		routes.GET("/:id/stops", getRouteStops)
// 	}

// 	// Bus stops
// 	stops := router.Group("/api/stops")
// 	{
// 		stops.GET("/", getAllStops)
// 		stops.GET("/:id", getStopById)
// 		stops.GET("/:id/arrivals", getStopArrivals)
// 	}

// 	// WebSocket endpoint
// 	router.GET("/ws", handleWebSocket)
// }

// // WebSocket handler
// func handleWebSocket(c *gin.Context) {
// 	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Error upgrading to WebSocket:", err)
// 		return
// 	}
// 	defer ws.Close()

// 	// Register client
// 	wsClients[ws] = true

// 	// Read messages from the client
// 	for {
// 		var msg Message
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			delete(wsClients, ws)
// 			break
// 		}

// 		// Handle the message based on its type
// 		switch msg.Type {
// 		case "location_update":
// 			// Process location update
// 			if update, ok := msg.Content.(map[string]interface{}); ok {
// 				// Use the update variable to access the location update data
// 				log.Println("Received location update:", update)
// 				// Convert to proper struct and process
// 				// Then broadcast to all clients
// 				broadcast <- msg
// 			}
// 		}
// 	}
// }

// func handleMessages() {
// 	for {
// 		// Get the next message from the broadcast channel
// 		msg := <-broadcast

// 		// Send it to all connected clients
// 		for client := range wsClients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				log.Println("Error sending message:", err)
// 				client.Close()
// 				delete(wsClients, client)
// 			}
// 		}
// 	}
// }

// // REST API Handlers

// func registerUser(c *gin.Context) {
// 	var user User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Hash password here in production

// 	result := db.Create(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"id": user.Model.ID, "username": user.Username, "role": user.Role})
// }

// func loginUser(c *gin.Context) {
// 	var credentials struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	if err := c.ShouldBindJSON(&credentials); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user User
// 	result := db.Where("username = ?", credentials.Username).First(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	// Check password here in production
// 	// For now, we'll skip password verification in this example

// 	// Generate JWT token here in production
// 	token := "sample-jwt-token"

// 	response := gin.H{
// 		"token": token,
// 		"user": gin.H{
// 			"id":       user.Model.ID,
// 			"username": user.Username,
// 			"role":     user.Role,
// 		},
// 	}

// 	// If user is a driver, include bus information
// 	if user.Role == "driver" {
// 		var bus Bus
// 		result := db.Where("driver_id = ?", user.ID).First(&bus)
// 		if result.Error == nil {
// 			response["bus"] = gin.H{
// 				"id":        bus.Model.ID,
// 				"busNumber": bus.BusNumber,
// 				"routeId":   bus.RouteID,
// 				"status":    bus.Status,
// 				"latitude":  bus.Latitude,
// 				"longitude": bus.Longitude,
// 			}
// 		} else {
// 			// No bus assigned
// 			response["bus"] = nil
// 		}
// 	}

// 	c.JSON(http.StatusOK, response)

// 	// Generate JWT token here in production
// 	// token := "sample-jwt-token"

// 	// c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.Model.ID, "username": user.Username, "role": user.Role}})
// }

// func getAllBuses(c *gin.Context) {
// 	var buses []Bus
// 	db.Find(&buses)
// 	c.JSON(http.StatusOK, buses)
// }

// func getBusById(c *gin.Context) {
// 	id := c.Param("id")
// 	var bus Bus

// 	result := db.First(&bus, id)
// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, bus)
// }

// func updateBusLocation(c *gin.Context) {
// 	id := c.Param("id")
// 	var locationUpdate struct {
// 		Latitude  float64 `json:"latitude"`
// 		Longitude float64 `json:"longitude"`
// 	}

// 	if err := c.ShouldBindJSON(&locationUpdate); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var bus Bus
// 	result := db.First(&bus, id)
// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
// 		return
// 	}

// 	// Update location
// 	bus.Latitude = locationUpdate.Latitude
// 	bus.Longitude = locationUpdate.Longitude
// 	bus.LastUpdated = time.Now()

// 	db.Save(&bus)

// 	// Broadcast update to all connected clients
// 	broadcast <- Message{
// 		Type: "location_update",
// 		Content: BusLocationUpdate{
// 			BusID:     bus.Model.ID,
// 			Latitude:  bus.Latitude,
// 			Longitude: bus.Longitude,
// 			Status:    bus.Status,
// 			Timestamp: time.Now().Unix(),
// 		},
// 	}

// 	c.JSON(http.StatusOK, bus)
// }

// func updateBusStatus(c *gin.Context) {
// 	id := c.Param("id")
// 	var statusUpdate struct {
// 		Status string `json:"status"`
// 	}

// 	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var bus Bus
// 	result := db.First(&bus, id)
// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
// 		return
// 	}

// 	// Update status
// 	bus.Status = statusUpdate.Status
// 	bus.LastUpdated = time.Now()

// 	db.Save(&bus)

// 	// Broadcast update
// 	broadcast <- Message{
// 		Type: "status_update",
// 		Content: gin.H{
// 			"busId":     bus.Model.ID,
// 			"status":    bus.Status,
// 			"timestamp": time.Now().Unix(),
// 		},
// 	}

// 	c.JSON(http.StatusOK, bus)
// }

// func getAllRoutes(c *gin.Context) {
// 	var routes []Route
// 	db.Find(&routes)
// 	c.JSON(http.StatusOK, routes)
// }

// func getRouteById(c *gin.Context) {
// 	id := c.Param("id")
// 	var route Route

// 	result := db.Preload("Stops", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("bus_stops.order ASC")
// 	}).First(&route, id)

// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, route)
// }

// func getRouteStops(c *gin.Context) {
// 	routeId := c.Param("id")
// 	var stops []BusStop

// 	result := db.Where("route_id = ?", routeId).Order("order ASC").Find(&stops)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stops"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, stops)
// }

// func getAllStops(c *gin.Context) {
// 	var stops []BusStop
// 	db.Find(&stops)
// 	c.JSON(http.StatusOK, stops)
// }

// func getStopById(c *gin.Context) {
// 	id := c.Param("id")
// 	var stop BusStop

// 	result := db.First(&stop, id)
// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, stop)
// }

// func getStopArrivals(c *gin.Context) {
// 	stopId := c.Param("id")

// 	var stop BusStop
// 	result := db.First(&stop, stopId)
// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
// 		return
// 	}
// 	// This would involve complex calculations in production
// 	// For now, return dummy data
// 	arrivals := []gin.H{
// 		{
// 			"busId":         1,
// 			"busNumber":     "101",
// 			"routeNumber":   "101",
// 			"estimatedTime": time.Now().Add(5 * time.Minute).Unix(),
// 		},
// 		{
// 			"busId":         2,
// 			"busNumber":     "202",
// 			"routeNumber":   "202",
// 			"estimatedTime": time.Now().Add(12 * time.Minute).Unix(),
// 		},
// 	}

// 	c.JSON(http.StatusOK, arrivals)
// }

// // HERE Maps API integration
// type HEREMapsService struct {
// 	ApiKey string
// }

// func NewHEREMapsService(apiKey string) *HEREMapsService {
// 	return &HEREMapsService{
// 		ApiKey: apiKey,
// 	}
// }

// func (h *HEREMapsService) CalculateETA(startLat, startLng, endLat, endLng float64) (int, error) {
// 	// In a real implementation, this would call the HERE Maps Routing API
// 	// For now, return a dummy value
// 	return 10, nil // 10 minutes
// }

// func (h *HEREMapsService) GetRoutePolyline(points []GeoPoint) ([]GeoPoint, error) {
// 	// In a real implementation, this would call the HERE Maps Routing API
// 	// For now, return the input points
// 	return points, nil
// }

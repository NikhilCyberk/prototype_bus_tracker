// // handlers/bus_handler.go
// package handlers

// import (
// 	"bustracking/models"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// // BusHandler handles bus-related requests
// type BusHandler struct {
// 	DB      *gorm.DB
// 	Emitter chan models.Message
// }

// // NewBusHandler creates a new bus handler
// func NewBusHandler(db *gorm.DB, emitter chan models.Message) *BusHandler {
// 	return &BusHandler{DB: db, Emitter: emitter}
// }

// // GetAllBuses returns all buses
// func (h *BusHandler) GetAllBuses(c *gin.Context) {
// 	var buses []models.Bus

// 	query := h.DB.Model(&models.Bus{})

// 	// Filter by route if provided
// 	if routeID := c.Query("routeId"); routeID != "" {
// 		query = query.Where("route_id = ?", routeID)
// 	}

// 	// Filter by status if provided
// 	if status := c.Query("status"); status != "" {
// 		query = query.Where("status = ?", status)
// 	}

// 	// Include related data
// 	query = query.Preload("Driver")
// 	query = query.Preload("Route")

// 	// Execute query
// 	if err := query.Find(&buses).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buses"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, buses)
// }

// // GetBusById returns a specific bus
// func (h *BusHandler) GetBusById(c *gin.Context) {
// 	id := c.Param("id")
// 	var bus models.Bus

// 	// Find the bus with its relationships
// 	if err := h.DB.Preload("Driver").Preload("Route").First(&bus, id).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, bus)
// }

// // UpdateBusLocation updates a bus's location
// func (h *BusHandler) UpdateBusLocation(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	var locationUpdate struct {
// 		Latitude  float64 `json:"latitude" binding:"required"`
// 		Longitude float64 `json:"longitude" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&locationUpdate); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if the bus exists
// 	var bus models.Bus
// 	if err := h.DB.First(&bus, id).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
// 		}
// 		return
// 	}

// 	// Authorize the update (only the assigned driver or an admin should be able to update)
// 	userID, _ := c.Get("userID")
// 	userRole, _ := c.Get("role")

// 	if userRole != "admin" && (userID.(uint) != bus.DriverID) {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this bus's location"})
// 		return
// 	}

// 	// Update location
// 	bus.Latitude = locationUpdate.Latitude
// 	bus.Longitude = locationUpdate.Longitude
// 	bus.LastUpdated = time.Now()

// 	if err := h.DB.Save(&bus).Error; err != nil {
// 		c.JSON(http.Status

// handlers/bus_handler.go
package handlers

import (
	"bustracking/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// BusHandler handles bus-related requests
type BusHandler struct {
	DB        *gorm.DB
	Broadcast chan models.Message
}

// NewBusHandler creates a new bus handler
func NewBusHandler(db *gorm.DB, broadcast chan models.Message) *BusHandler {
	return &BusHandler{
		DB:        db,
		Broadcast: broadcast,
	}
}

// GetAllBuses returns all buses with optional filters
func (h *BusHandler) GetAllBuses(c *gin.Context) {
	var buses []models.Bus

	query := h.DB.Model(&models.Bus{}).Preload("Driver").Preload("Route")

	// Apply filters if provided
	if routeID := c.Query("routeId"); routeID != "" {
		query = query.Where("route_id = ?", routeID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if busNumber := c.Query("busNumber"); busNumber != "" {
		query = query.Where("bus_number LIKE ?", "%"+busNumber+"%")
	}

	// Execute query
	if err := query.Find(&buses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buses"})
		return
	}

	c.JSON(http.StatusOK, buses)
}

// GetBusById returns a specific bus by ID
func (h *BusHandler) GetBusById(c *gin.Context) {
	id := c.Param("id")
	var bus models.Bus

	if err := h.DB.Preload("Driver").Preload("Route").First(&bus, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
		}
		return
	}

	c.JSON(http.StatusOK, bus)
}

// UpdateBusLocation updates a bus's location
func (h *BusHandler) UpdateBusLocation(c *gin.Context) {
	id := c.Param("id")
	var locationUpdate struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&locationUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bus models.Bus
	if err := h.DB.First(&bus, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
		}
		return
	}

	// Check if the current user is the driver of this bus or an admin
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if userRole != "admin" && userID.(uint) != bus.DriverID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this bus's location"})
		return
	}

	// Update location
	bus.Latitude = locationUpdate.Latitude
	bus.Longitude = locationUpdate.Longitude
	bus.LastUpdated = time.Now()

	if err := h.DB.Save(&bus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bus location"})
		return
	}

	// Broadcast the update to all WebSocket clients
	h.Broadcast <- models.Message{
		Type: "location_update",
		Content: models.BusLocationUpdate{
			BusID:     bus.ID,
			Latitude:  bus.Latitude,
			Longitude: bus.Longitude,
			Status:    bus.Status,
			Timestamp: time.Now().Unix(),
		},
	}

	c.JSON(http.StatusOK, bus)
}

// UpdateBusStatus updates a bus's status
func (h *BusHandler) UpdateBusStatus(c *gin.Context) {
	id := c.Param("id")
	var statusUpdate struct {
		Status string `json:"status" binding:"required,oneof=on-route delayed off-duty"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bus models.Bus
	if err := h.DB.First(&bus, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
		}
		return
	}

	// Check if the current user is the driver of this bus or an admin
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if userRole != "admin" && userID.(uint) != bus.DriverID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this bus's status"})
		return
	}

	// Update status
	bus.Status = statusUpdate.Status
	bus.LastUpdated = time.Now()

	if err := h.DB.Save(&bus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bus status"})
		return
	}

	// Broadcast the update to all WebSocket clients
	h.Broadcast <- models.Message{
		Type: "status_update",
		Content: gin.H{
			"busId":     bus.ID,
			"status":    bus.Status,
			"timestamp": time.Now().Unix(),
		},
	}

	c.JSON(http.StatusOK, bus)
}

// AssignBusToDriver assigns a bus to a driver (admin only)
func (h *BusHandler) AssignBusToDriver(c *gin.Context) {
	var assignment struct {
		DriverID uint `json:"driverId" binding:"required"`
		BusID    uint `json:"busId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if driver exists and is actually a driver
	var driver models.User
	if err := h.DB.First(&driver, assignment.DriverID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch driver"})
		}
		return
	}

	if driver.Role != "driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a driver"})
		return
	}

	// Check if bus exists
	var bus models.Bus
	if err := h.DB.First(&bus, assignment.BusID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
		}
		return
	}

	// Check if the bus is already assigned to another driver
	if bus.DriverID != 0 && bus.DriverID != assignment.DriverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bus is already assigned to another driver"})
		return
	}

	// Update bus with new driver
	bus.DriverID = assignment.DriverID
	if err := h.DB.Save(&bus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign bus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bus assigned successfully",
		"bus": gin.H{
			"id":        bus.ID,
			"busNumber": bus.BusNumber,
			"driverId":  bus.DriverID,
			"routeId":   bus.RouteID,
		},
	})
}

// CreateBus creates a new bus (admin only)
func (h *BusHandler) CreateBus(c *gin.Context) {
	var bus struct {
		BusNumber string  `json:"busNumber" binding:"required"`
		RouteID   uint    `json:"routeId" binding:"required"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if route exists
	var route models.Route
	if err := h.DB.First(&route, bus.RouteID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch route"})
		}
		return
	}

	newBus := models.Bus{
		BusNumber:   bus.BusNumber,
		RouteID:     bus.RouteID,
		Status:      "off-duty",
		Latitude:    bus.Latitude,
		Longitude:   bus.Longitude,
		LastUpdated: time.Now(),
	}

	if err := h.DB.Create(&newBus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bus"})
		return
	}

	c.JSON(http.StatusCreated, newBus)
}

// UpdateBus updates bus information (admin only)
func (h *BusHandler) UpdateBus(c *gin.Context) {
	id := c.Param("id")
	var update struct {
		BusNumber string  `json:"busNumber"`
		RouteID   uint    `json:"routeId"`
		DriverID  uint    `json:"driverId"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Status    string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bus models.Bus
	if err := h.DB.First(&bus, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
		}
		return
	}

	// Update fields if they are provided
	if update.BusNumber != "" {
		bus.BusNumber = update.BusNumber
	}
	if update.RouteID != 0 {
		bus.RouteID = update.RouteID
	}
	if update.DriverID != 0 {
		bus.DriverID = update.DriverID
	}
	if update.Latitude != 0 {
		bus.Latitude = update.Latitude
	}
	if update.Longitude != 0 {
		bus.Longitude = update.Longitude
	}
	if update.Status != "" {
		bus.Status = update.Status
	}

	bus.LastUpdated = time.Now()

	if err := h.DB.Save(&bus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bus"})
		return
	}

	c.JSON(http.StatusOK, bus)
}

// DeleteBus deletes a bus (admin only)
func (h *BusHandler) DeleteBus(c *gin.Context) {
	id := c.Param("id")

	var bus models.Bus
	if err := h.DB.First(&bus, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bus"})
		}
		return
	}

	if err := h.DB.Delete(&bus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bus deleted successfully"})
}

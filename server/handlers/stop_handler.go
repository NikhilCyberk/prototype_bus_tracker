// handlers/stop_handler.go
package handlers

import (
	"bustracking/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StopHandler struct {
	DB *gorm.DB
}

func NewStopHandler(db *gorm.DB) *StopHandler {
	return &StopHandler{DB: db}
}

func (h *StopHandler) GetAllStops(c *gin.Context) {
	var stops []models.BusStop

	if err := h.DB.Find(&stops).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stops"})
		return
	}

	c.JSON(http.StatusOK, stops)
}

func (h *StopHandler) GetStopById(c *gin.Context) {
	id := c.Param("id")
	var stop models.BusStop

	if err := h.DB.First(&stop, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stop"})
		}
		return
	}

	c.JSON(http.StatusOK, stop)
}

func (h *StopHandler) GetStopArrivals(c *gin.Context) {
	stopID := c.Param("id")
	var stop models.BusStop

	// First get the stop details
	if err := h.DB.First(&stop, stopID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stop"})
		}
		return
	}

	// Get all buses on this route
	var buses []models.Bus
	if err := h.DB.Where("route_id = ?", stop.RouteID).Find(&buses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buses"})
		return
	}

	// In a real implementation, you would calculate ETAs here
	// For now, we'll return dummy data
	arrivals := make([]models.StopArrival, len(buses))
	for i, bus := range buses {
		arrivals[i] = models.StopArrival{
			BusID:         bus.ID,
			BusNumber:     bus.BusNumber,
			RouteNumber:   bus.Route.RouteNumber,
			EstimatedTime: int64((i + 1) * 5 * 60), // Dummy ETA in seconds (5, 10, 15... minutes)
		}
	}

	c.JSON(http.StatusOK, arrivals)
}

func (h *StopHandler) CreateStop(c *gin.Context) {
	var stop models.BusStop
	if err := c.ShouldBindJSON(&stop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&stop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stop"})
		return
	}

	c.JSON(http.StatusCreated, stop)
}

func (h *StopHandler) UpdateStop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var stop models.BusStop
	if err := h.DB.First(&stop, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stop"})
		}
		return
	}

	if err := c.ShouldBindJSON(&stop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&stop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stop"})
		return
	}

	c.JSON(http.StatusOK, stop)
}

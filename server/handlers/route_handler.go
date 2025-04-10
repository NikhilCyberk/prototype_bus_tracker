// handlers/route_handler.go
package handlers

import (
	"bustracking/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouteHandler struct {
	DB *gorm.DB
}

func NewRouteHandler(db *gorm.DB) *RouteHandler {
	return &RouteHandler{DB: db}
}

func (h *RouteHandler) GetAllRoutes(c *gin.Context) {
	var routes []models.Route

	// Preload stops and order them
	if err := h.DB.Preload("Stops", func(db *gorm.DB) *gorm.DB {
		return db.Order("bus_stops.order ASC")
	}).Find(&routes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routes"})
		return
	}

	c.JSON(http.StatusOK, routes)
}

func (h *RouteHandler) GetRouteById(c *gin.Context) {
	id := c.Param("id")
	var route models.Route

	if err := h.DB.Preload("Stops", func(db *gorm.DB) *gorm.DB {
		return db.Order("bus_stops.order ASC")
	}).First(&route, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch route"})
		}
		return
	}

	c.JSON(http.StatusOK, route)
}

func (h *RouteHandler) GetRouteStops(c *gin.Context) {
	routeID := c.Param("id")
	var stops []models.BusStop

	if err := h.DB.Where("route_id = ?", routeID).Order("order ASC").Find(&stops).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stops"})
		return
	}

	c.JSON(http.StatusOK, stops)
}

func (h *RouteHandler) CreateRoute(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create route"})
		return
	}

	c.JSON(http.StatusCreated, route)
}

func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var route models.Route
	if err := h.DB.First(&route, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch route"})
		}
		return
	}

	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route"})
		return
	}

	c.JSON(http.StatusOK, route)
}

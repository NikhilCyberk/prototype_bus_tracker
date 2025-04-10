// models/bus.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Bus represents a bus in the system
type Bus struct {
	gorm.Model
	BusNumber   string    `json:"busNumber"`
	DriverID    uint      `json:"driverId"`
	Driver      User      `json:"driver" gorm:"foreignKey:DriverID"`
	RouteID     uint      `json:"routeId"`
	Route       Route     `json:"route" gorm:"foreignKey:RouteID"`
	Status      string    `json:"status"` // "on-route", "delayed", "off-duty"
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// TableName returns the table name for the bus model
func (Bus) TableName() string {
	return "buses"
}

// BusLocationUpdate represents a location update for a bus
type BusLocationUpdate struct {
	BusID     uint    `json:"busId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
	Timestamp int64   `json:"timestamp"`
}

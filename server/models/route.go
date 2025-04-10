// models/route.go
package models

import (
	"gorm.io/gorm"
)

// Route represents a bus route
type Route struct {
	gorm.Model
	RouteNumber string     `json:"routeNumber"`
	Name        string     `json:"name"`
	Color       string     `json:"color"`
	Stops       []BusStop  `json:"stops" gorm:"foreignKey:RouteID"`
	Path        []GeoPoint `gorm:"serializer:json" json:"path"`
}

// GeoPoint represents a geographic coordinate
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TableName returns the table name for the route model
func (Route) TableName() string {
	return "routes"
}

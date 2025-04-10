// models/stop.go
package models

import (
	"gorm.io/gorm"
)

// BusStop represents a stop on a bus route
type BusStop struct {
	gorm.Model
	Name      string  `json:"name"`
	RouteID   uint    `json:"routeId"`
	Order     int     `json:"order"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TableName returns the table name for the bus stop model
func (BusStop) TableName() string {
	return "bus_stops"
}

// StopArrival represents an estimated arrival of a bus at a stop
type StopArrival struct {
	BusID         uint   `json:"busId"`
	BusNumber     string `json:"busNumber"`
	RouteNumber   string `json:"routeNumber"`
	EstimatedTime int64  `json:"estimatedTime"`
}

// services/maps_service.go
package services

import (
	"bustracking/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// MapsService provides map-related functionality
type MapsService struct {
	ApiKey     string
	HttpClient *http.Client
}

// NewMapsService creates a new maps service instance
func NewMapsService(apiKey string) *MapsService {
	return &MapsService{
		ApiKey: apiKey,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// CalculateETA estimates time of arrival between two points
func (m *MapsService) CalculateETA(startLat, startLng, endLat, endLng float64) (int, error) {
	if m.ApiKey == "" {
		// Return a simulated ETA if no API key (for demo purposes)
		return simulateETA(startLat, startLng, endLat, endLng), nil
	}

	// Construct the URL for the HERE Maps Routing API
	baseURL := "https://router.hereapi.com/v8/routes"
	params := url.Values{}
	params.Add("apiKey", m.ApiKey)
	params.Add("transportMode", "bus")
	params.Add("origin", fmt.Sprintf("%f,%f", startLat, startLng))
	params.Add("destination", fmt.Sprintf("%f,%f", endLat, endLng))
	params.Add("return", "summary")

	// Make the API request
	resp, err := m.HttpClient.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("maps API error: status code %d", resp.StatusCode)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Extract the duration from the response
	routes, ok := result["routes"].([]interface{})
	if !ok || len(routes) == 0 {
		return 0, errors.New("no routes found")
	}

	route, ok := routes[0].(map[string]interface{})
	if !ok {
		return 0, errors.New("invalid route data")
	}

	sections, ok := route["sections"].([]interface{})
	if !ok || len(sections) == 0 {
		return 0, errors.New("no sections found")
	}

	section, ok := sections[0].(map[string]interface{})
	if !ok {
		return 0, errors.New("invalid section data")
	}

	summary, ok := section["summary"].(map[string]interface{})
	if !ok {
		return 0, errors.New("no summary found")
	}

	duration, ok := summary["duration"].(float64)
	if !ok {
		return 0, errors.New("invalid duration data")
	}

	// Convert duration from seconds to minutes
	return int(duration / 60), nil
}

// GetRoutePolyline generates a polyline for a route
func (m *MapsService) GetRoutePolyline(points []models.GeoPoint) ([]models.GeoPoint, error) {
	if m.ApiKey == "" || len(points) < 2 {
		// Return the input points if no API key or not enough points
		return points, nil
	}

	// Build waypoints string
	var waypoints string
	for i, point := range points {
		if i > 0 {
			waypoints += ";"
		}
		waypoints += fmt.Sprintf("%f,%f", point.Latitude, point.Longitude)
	}

	// Construct the URL for the HERE Maps Routing API
	baseURL := "https://router.hereapi.com/v8/routes"
	params := url.Values{}
	params.Add("apiKey", m.ApiKey)
	params.Add("transportMode", "bus")
	params.Add("return", "polyline")
	params.Add("routingMode", "short")
	params.Add("origin", fmt.Sprintf("%f,%f", points[0].Latitude, points[0].Longitude))
	params.Add("destination", fmt.Sprintf("%f,%f", points[len(points)-1].Latitude, points[len(points)-1].Longitude))

	if len(points) > 2 {
		var viaPoints string
		for i := 1; i < len(points)-1; i++ {
			if i > 1 {
				viaPoints += ";"
			}
			viaPoints += fmt.Sprintf("%f,%f", points[i].Latitude, points[i].Longitude)
		}
		params.Add("via", viaPoints)
	}

	// Make the API request
	resp, err := m.HttpClient.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return points, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return points, fmt.Errorf("maps API error: status code %d", resp.StatusCode)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return points, err
	}

	// Extract the polyline from the response
	routes, ok := result["routes"].([]interface{})
	if !ok || len(routes) == 0 {
		return points, errors.New("no routes found")
	}

	// Process polyline and convert to GeoPoint array
	// (This would be more complex in real implementation)

	// For now, just return the original points
	return points, nil
}

// simulateETA provides a simple ETA calculation for demo purposes
func simulateETA(startLat, startLng, endLat, endLng float64) int {
	// Calculate rough distance using simplified formula
	latDiff := startLat - endLat
	lngDiff := startLng - endLng

	// Simple Euclidean distance (not accurate for real world)
	distance := (latDiff*latDiff + lngDiff*lngDiff) * 111.32 * 1000 // Rough meters

	// Assume average bus speed of 20 km/h = 5.55 m/s
	timeInSeconds := distance / 5.55

	// Convert to minutes and add some variability
	minutes := int(timeInSeconds / 60)

	// Ensure minimum 1 minute, maximum 60 minutes for demo
	if minutes < 1 {
		minutes = 1
	} else if minutes > 60 {
		minutes = 60
	}

	return minutes
}

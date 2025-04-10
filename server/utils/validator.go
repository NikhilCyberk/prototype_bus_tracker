// utils/validator.go
package utils

import (
	"regexp"
	"strings"
)

func ValidateUsername(username string) bool {
	// Username must be 3-20 characters long, alphanumeric with underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
	return matched
}

func ValidatePassword(password string) bool {
	// Password must be at least 8 characters long
	return len(password) >= 8
}

func ValidateRole(role string) bool {
	role = strings.ToLower(role)
	return role == "driver" || role == "passenger" || role == "admin"
}

func ValidateCoordinates(lat, lng float64) bool {
	// Basic validation for latitude and longitude
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

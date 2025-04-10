// models/user.go
package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`    // Excludes password from JSON
	Role     string `json:"role"` // "driver" or "passenger" or "admin"
}

// TableName returns the table name for the user model
func (User) TableName() string {
	return "users"
}

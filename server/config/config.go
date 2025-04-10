// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration stores app settings
type Configuration struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	ServerPort string
	JWTSecret  string
	MapsAPIKey string
}

// LoadConfig loads configuration from environment variables or uses defaults
func LoadConfig() *Configuration {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config := &Configuration{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "1234"),
		DBName:     getEnv("DB_NAME", "bustrackerdb"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		MapsAPIKey: getEnv("MAPS_API_KEY", ""),
	}

	return config
}

// GetDSN returns the database connection string
func (c *Configuration) GetDSN() string {
	return "host=" + c.DBHost + " user=" + c.DBUser +
		" password=" + c.DBPassword + " dbname=" + c.DBName +
		" port=" + c.DBPort + " sslmode=disable"
}

// GetServerAddress returns the server address with port
func (c *Configuration) GetServerAddress() string {
	return ":" + c.ServerPort
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if fallback == "" && (key == "JWT_SECRET") {
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	return fallback
}

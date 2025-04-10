// seeder/seeder.go
package seeder

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Database models (duplicated from main.go for this example)
// In a real project, you would import these from a models package
type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"` // "driver" or "passenger"
}

type Bus struct {
	gorm.Model
	BusNumber   string    `json:"busNumber"`
	DriverID    uint      `json:"driverId"`
	Driver      User      `json:"driver"`
	RouteID     uint      `json:"routeId"`
	Route       Route     `json:"route"`
	Status      string    `json:"status"` // "on-route", "delayed", "off-duty"
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type Route struct {
	gorm.Model
	RouteNumber string     `json:"routeNumber"`
	Name        string     `json:"name"`
	Color       string     `json:"color"`
	Stops       []BusStop  `json:"stops"`
	Path        []GeoPoint `gorm:"serializer:json" json:"path"`
}

type BusStop struct {
	gorm.Model
	Name      string  `json:"name"`
	RouteID   uint    `json:"routeId"`
	Order     int     `json:"order"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Seeder struct to hold the database connection
type Seeder struct {
	DB *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{DB: db}
}

// hashPassword generates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// // SeedIndianData creates sample data with Indian context
// func (s *Seeder) SeedIndianData() error {
// 	// Check if we need to seed data
// 	var count int64
// 	s.DB.Model(&Route{}).Count(&count)
// 	if count > 0 {
// 		log.Println("Data already exists, skipping seeding")
// 		return nil // Data already exists
// 	}

// 	// Seed users - both passengers and drivers
// 	users := []User{
// 		// Passengers
// 		{
// 			Username: "raj_kumar",
// 			Password: "password123", // In production, this would be hashed
// 			Role:     "passenger",
// 		},
// 		{
// 			Username: "priya_sharma",
// 			Password: "password123",
// 			Role:     "passenger",
// 		},
// 		{
// 			Username: "amit_patel",
// 			Password: "password123",
// 			Role:     "passenger",
// 		},
// 		// Drivers
// 		{
// 			Username: "suresh_driver",
// 			Password: "driver123",
// 			Role:     "driver",
// 		},
// 		{
// 			Username: "vikram_singh",
// 			Password: "driver123",
// 			Role:     "driver",
// 		},
// 		{
// 			Username: "deepak_verma",
// 			Password: "driver123",
// 			Role:     "driver",
// 		},
// 	}

// 	log.Println("Creating users...")
// 	for i := range users {
// 		hashedPassword, err := hashPassword(users[i].Password)
// 		if err != nil {
// 			return err
// 		}
// 		users[i].Password = hashedPassword
// 		if err := s.DB.Create(&users[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// Seed routes for Delhi
// 	routes := []Route{
// 		{
// 			RouteNumber: "DL101",
// 			Name:        "Connaught Place to India Gate",
// 			Color:       "#FF0000",
// 			Path: []GeoPoint{
// 				{Latitude: 28.6289, Longitude: 77.2311}, // CP
// 				{Latitude: 28.6270, Longitude: 77.2330},
// 				{Latitude: 28.6250, Longitude: 77.2350},
// 				{Latitude: 28.6129, Longitude: 77.2295}, // India Gate
// 			},
// 		},
// 		{
// 			RouteNumber: "DL202",
// 			Name:        "Karol Bagh to Red Fort",
// 			Color:       "#00FF00",
// 			Path: []GeoPoint{
// 				{Latitude: 28.6449, Longitude: 77.1905}, // Karol Bagh
// 				{Latitude: 28.6460, Longitude: 77.2000},
// 				{Latitude: 28.6471, Longitude: 77.2100},
// 				{Latitude: 28.6526, Longitude: 77.2311}, // Red Fort
// 			},
// 		},
// 		{
// 			RouteNumber: "DL303",
// 			Name:        "Lajpat Nagar to Chandni Chowk",
// 			Color:       "#0000FF",
// 			Path: []GeoPoint{
// 				{Latitude: 28.5700, Longitude: 77.2400}, // Lajpat Nagar
// 				{Latitude: 28.5900, Longitude: 77.2350},
// 				{Latitude: 28.6100, Longitude: 77.2300},
// 				{Latitude: 28.6505, Longitude: 77.2303}, // Chandni Chowk
// 			},
// 		},
// 	}

// 	log.Println("Creating routes...")
// 	// Create routes
// 	for i := range routes {
// 		if err := s.DB.Create(&routes[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// Create stops for each route
// 	log.Println("Creating bus stops...")

// 	// DL101 stops
// 	dl101Stops := []BusStop{
// 		{
// 			Name:      "Connaught Place",
// 			RouteID:   routes[0].ID,
// 			Order:     1,
// 			Latitude:  28.6289,
// 			Longitude: 77.2311,
// 		},
// 		{
// 			Name:      "Janpath",
// 			RouteID:   routes[0].ID,
// 			Order:     2,
// 			Latitude:  28.6270,
// 			Longitude: 77.2330,
// 		},
// 		{
// 			Name:      "Windsor Place",
// 			RouteID:   routes[0].ID,
// 			Order:     3,
// 			Latitude:  28.6250,
// 			Longitude: 77.2350,
// 		},
// 		{
// 			Name:      "India Gate",
// 			RouteID:   routes[0].ID,
// 			Order:     4,
// 			Latitude:  28.6129,
// 			Longitude: 77.2295,
// 		},
// 	}

// 	for i := range dl101Stops {
// 		if err := s.DB.Create(&dl101Stops[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// DL202 stops
// 	dl202Stops := []BusStop{
// 		{
// 			Name:      "Karol Bagh",
// 			RouteID:   routes[1].ID,
// 			Order:     1,
// 			Latitude:  28.6449,
// 			Longitude: 77.1905,
// 		},
// 		{
// 			Name:      "Patel Nagar",
// 			RouteID:   routes[1].ID,
// 			Order:     2,
// 			Latitude:  28.6460,
// 			Longitude: 77.2000,
// 		},
// 		{
// 			Name:      "New Delhi Railway Station",
// 			RouteID:   routes[1].ID,
// 			Order:     3,
// 			Latitude:  28.6471,
// 			Longitude: 77.2100,
// 		},
// 		{
// 			Name:      "Red Fort",
// 			RouteID:   routes[1].ID,
// 			Order:     4,
// 			Latitude:  28.6526,
// 			Longitude: 77.2311,
// 		},
// 	}

// 	for i := range dl202Stops {
// 		if err := s.DB.Create(&dl202Stops[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// DL303 stops
// 	dl303Stops := []BusStop{
// 		{
// 			Name:      "Lajpat Nagar",
// 			RouteID:   routes[2].ID,
// 			Order:     1,
// 			Latitude:  28.5700,
// 			Longitude: 77.2400,
// 		},
// 		{
// 			Name:      "Moolchand",
// 			RouteID:   routes[2].ID,
// 			Order:     2,
// 			Latitude:  28.5900,
// 			Longitude: 77.2350,
// 		},
// 		{
// 			Name:      "ITO",
// 			RouteID:   routes[2].ID,
// 			Order:     3,
// 			Latitude:  28.6100,
// 			Longitude: 77.2300,
// 		},
// 		{
// 			Name:      "Chandni Chowk",
// 			RouteID:   routes[2].ID,
// 			Order:     4,
// 			Latitude:  28.6505,
// 			Longitude: 77.2303,
// 		},
// 	}

// 	for i := range dl303Stops {
// 		if err := s.DB.Create(&dl303Stops[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// Get driver IDs from database to ensure we have the correct IDs
// 	var drivers []User
// 	if err := s.DB.Where("role = ?", "driver").Find(&drivers).Error; err != nil {
// 		return err
// 	}

// 	log.Println("Creating buses...")
// 	// Create buses with assigned drivers
// 	buses := []Bus{
// 		{
// 			BusNumber:   "DL-1S-1234",
// 			DriverID:    drivers[0].ID,
// 			RouteID:     routes[0].ID,
// 			Status:      "on-route",
// 			Latitude:    28.6289,
// 			Longitude:   77.2311,
// 			LastUpdated: time.Now(),
// 		},
// 		{
// 			BusNumber:   "DL-2V-5678",
// 			DriverID:    drivers[1].ID,
// 			RouteID:     routes[1].ID,
// 			Status:      "on-route",
// 			Latitude:    28.6449,
// 			Longitude:   77.1905,
// 			LastUpdated: time.Now(),
// 		},
// 		{
// 			BusNumber:   "DL-3D-9012",
// 			DriverID:    drivers[2].ID,
// 			RouteID:     routes[2].ID,
// 			Status:      "on-route",
// 			Latitude:    28.5700,
// 			Longitude:   77.2400,
// 			LastUpdated: time.Now(),
// 		},
// 	}

// 	for i := range buses {
// 		if err := s.DB.Create(&buses[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	log.Println("Seed data with Indian context created successfully")
// 	return nil
// }

func (s *Seeder) SeedIndianData() error {
	// Check if we need to seed data
	var count int64
	s.DB.Model(&Route{}).Count(&count)
	if count > 0 {
		log.Println("Data already exists, skipping seeding")
		return nil // Data already exists
	}

	// Seed users - both passengers and drivers
	users := []User{
		// Passengers
		{
			Username: "rahul_sharma",
			Password: "password123",
			Role:     "passenger",
		},
		{
			Username: "priya_patel",
			Password: "password123",
			Role:     "passenger",
		},
		{
			Username: "amit_kumar",
			Password: "password123",
			Role:     "passenger",
		},
		{
			Username: "neha_verma",
			Password: "password123",
			Role:     "passenger",
		},
		{
			Username: "arjun_singh",
			Password: "password123",
			Role:     "passenger",
		},
		// Drivers
		{
			Username: "suresh_reddy",
			Password: "driver123",
			Role:     "driver",
		},
		{
			Username: "vikram_naidu",
			Password: "driver123",
			Role:     "driver",
		},
		{
			Username: "deepak_mishra",
			Password: "driver123",
			Role:     "driver",
		},
		{
			Username: "rajesh_kumar",
			Password: "driver123",
			Role:     "driver",
		},
		{
			Username: "mohan_sharma",
			Password: "driver123",
			Role:     "driver",
		},
	}

	log.Println("Creating users...")
	for i := range users {
		hashedPassword, err := hashPassword(users[i].Password)
		if err != nil {
			return err
		}
		users[i].Password = hashedPassword
		if err := s.DB.Create(&users[i]).Error; err != nil {
			return err
		}
	}

	// Seed routes for major Indian cities
	routes := []Route{
		// Delhi Routes
		{
			RouteNumber: "DL-101",
			Name:        "Kashmere Gate to Qutub Minar",
			Color:       "#FF5733",
			Path: []GeoPoint{
				{Latitude: 28.6654, Longitude: 77.2311}, // Kashmere Gate
				{Latitude: 28.6449, Longitude: 77.2183}, // Civil Lines
				{Latitude: 28.6289, Longitude: 77.2067}, // Connaught Place
				{Latitude: 28.6129, Longitude: 77.2295}, // India Gate
				{Latitude: 28.5929, Longitude: 77.2193}, // AIIMS
				{Latitude: 28.5245, Longitude: 77.1855}, // Qutub Minar
			},
		},
		{
			RouteNumber: "DL-202",
			Name:        "Anand Vihar to Dwarka",
			Color:       "#33FF57",
			Path: []GeoPoint{
				{Latitude: 28.6504, Longitude: 77.3162}, // Anand Vihar
				{Latitude: 28.6449, Longitude: 77.2532}, // ITO
				{Latitude: 28.6289, Longitude: 77.2067}, // Connaught Place
				{Latitude: 28.6129, Longitude: 77.2295}, // India Gate
				{Latitude: 28.5800, Longitude: 77.2000}, // Dhaula Kuan
				{Latitude: 28.5923, Longitude: 77.0419}, // Dwarka Sec 14
			},
		},

		// Mumbai Routes
		{
			RouteNumber: "MB-501",
			Name:        "CST to Borivali",
			Color:       "#3357FF",
			Path: []GeoPoint{
				{Latitude: 18.9398, Longitude: 72.8355}, // Chhatrapati Shivaji Terminus
				{Latitude: 18.9250, Longitude: 72.8345}, // Churchgate
				{Latitude: 19.0760, Longitude: 72.8777}, // Bandra
				{Latitude: 19.1136, Longitude: 72.8697}, // Andheri
				{Latitude: 19.2144, Longitude: 72.8479}, // Borivali
			},
		},

		// Bangalore Routes
		{
			RouteNumber: "BL-301",
			Name:        "Majestic to Electronic City",
			Color:       "#F033FF",
			Path: []GeoPoint{
				{Latitude: 12.9774, Longitude: 77.5661}, // Kempegowda Bus Station
				{Latitude: 12.9716, Longitude: 77.5946}, // Shanthinagar
				{Latitude: 12.9352, Longitude: 77.6245}, // Silk Board
				{Latitude: 12.8456, Longitude: 77.6483}, // Electronic City
			},
		},

		// Chennai Routes
		{
			RouteNumber: "CH-401",
			Name:        "CMBT to Thiruvanmiyur",
			Color:       "#FF33F0",
			Path: []GeoPoint{
				{Latitude: 13.0827, Longitude: 80.2077}, // CMBT
				{Latitude: 13.0629, Longitude: 80.2239}, // T Nagar
				{Latitude: 13.0399, Longitude: 80.2405}, // Saidapet
				{Latitude: 12.9855, Longitude: 80.2581}, // Thiruvanmiyur
			},
		},
	}

	log.Println("Creating routes...")
	// Create routes
	for i := range routes {
		if err := s.DB.Create(&routes[i]).Error; err != nil {
			return err
		}
	}

	// Create stops for each route
	log.Println("Creating bus stops...")

	// DL-101 stops (Kashmere Gate to Qutub Minar)
	dl101Stops := []BusStop{
		{Name: "Kashmere Gate ISBT", RouteID: routes[0].ID, Order: 1, Latitude: 28.6654, Longitude: 77.2311},
		{Name: "Civil Lines", RouteID: routes[0].ID, Order: 2, Latitude: 28.6449, Longitude: 77.2183},
		{Name: "Connaught Place", RouteID: routes[0].ID, Order: 3, Latitude: 28.6289, Longitude: 77.2067},
		{Name: "India Gate", RouteID: routes[0].ID, Order: 4, Latitude: 28.6129, Longitude: 77.2295},
		{Name: "AIIMS", RouteID: routes[0].ID, Order: 5, Latitude: 28.5929, Longitude: 77.2193},
		{Name: "Qutub Minar", RouteID: routes[0].ID, Order: 6, Latitude: 28.5245, Longitude: 77.1855},
	}

	// DL-202 stops (Anand Vihar to Dwarka)
	dl202Stops := []BusStop{
		{Name: "Anand Vihar ISBT", RouteID: routes[1].ID, Order: 1, Latitude: 28.6504, Longitude: 77.3162},
		{Name: "ITO", RouteID: routes[1].ID, Order: 2, Latitude: 28.6449, Longitude: 77.2532},
		{Name: "Connaught Place", RouteID: routes[1].ID, Order: 3, Latitude: 28.6289, Longitude: 77.2067},
		{Name: "India Gate", RouteID: routes[1].ID, Order: 4, Latitude: 28.6129, Longitude: 77.2295},
		{Name: "Dhaula Kuan", RouteID: routes[1].ID, Order: 5, Latitude: 28.5800, Longitude: 77.2000},
		{Name: "Dwarka Sector 14", RouteID: routes[1].ID, Order: 6, Latitude: 28.5923, Longitude: 77.0419},
	}

	// MB-501 stops (CST to Borivali)
	mb501Stops := []BusStop{
		{Name: "Chhatrapati Shivaji Terminus", RouteID: routes[2].ID, Order: 1, Latitude: 18.9398, Longitude: 72.8355},
		{Name: "Churchgate", RouteID: routes[2].ID, Order: 2, Latitude: 18.9250, Longitude: 72.8345},
		{Name: "Bandra", RouteID: routes[2].ID, Order: 3, Latitude: 19.0760, Longitude: 72.8777},
		{Name: "Andheri", RouteID: routes[2].ID, Order: 4, Latitude: 19.1136, Longitude: 72.8697},
		{Name: "Borivali", RouteID: routes[2].ID, Order: 5, Latitude: 19.2144, Longitude: 72.8479},
	}

	// BL-301 stops (Majestic to Electronic City)
	bl301Stops := []BusStop{
		{Name: "Kempegowda Bus Station", RouteID: routes[3].ID, Order: 1, Latitude: 12.9774, Longitude: 77.5661},
		{Name: "Shanthinagar", RouteID: routes[3].ID, Order: 2, Latitude: 12.9716, Longitude: 77.5946},
		{Name: "Silk Board", RouteID: routes[3].ID, Order: 3, Latitude: 12.9352, Longitude: 77.6245},
		{Name: "Electronic City", RouteID: routes[3].ID, Order: 4, Latitude: 12.8456, Longitude: 77.6483},
	}

	// CH-401 stops (CMBT to Thiruvanmiyur)
	ch401Stops := []BusStop{
		{Name: "CMBT", RouteID: routes[4].ID, Order: 1, Latitude: 13.0827, Longitude: 80.2077},
		{Name: "T Nagar", RouteID: routes[4].ID, Order: 2, Latitude: 13.0629, Longitude: 80.2239},
		{Name: "Saidapet", RouteID: routes[4].ID, Order: 3, Latitude: 13.0399, Longitude: 80.2405},
		{Name: "Thiruvanmiyur", RouteID: routes[4].ID, Order: 4, Latitude: 12.9855, Longitude: 80.2581},
	}

	// Combine all stops and create them
	allStops := append(dl101Stops, dl202Stops...)
	allStops = append(allStops, mb501Stops...)
	allStops = append(allStops, bl301Stops...)
	allStops = append(allStops, ch401Stops...)

	for i := range allStops {
		if err := s.DB.Create(&allStops[i]).Error; err != nil {
			return err
		}
	}

	// Get driver IDs from database
	var drivers []User
	if err := s.DB.Where("role = ?", "driver").Find(&drivers).Error; err != nil {
		return err
	}

	log.Println("Creating buses...")
	// Create buses with assigned drivers
	buses := []Bus{
		// Delhi Buses
		{
			BusNumber:   "DL-1PA-2345",
			DriverID:    drivers[0].ID,
			RouteID:     routes[0].ID,
			Status:      "on-route",
			Latitude:    28.6289,
			Longitude:   77.2067,
			LastUpdated: time.Now(),
		},
		{
			BusNumber:   "DL-1PA-6789",
			DriverID:    drivers[1].ID,
			RouteID:     routes[1].ID,
			Status:      "on-route",
			Latitude:    28.6449,
			Longitude:   77.2532,
			LastUpdated: time.Now(),
		},

		// Mumbai Buses
		{
			BusNumber:   "MH-02-AB-1234",
			DriverID:    drivers[2].ID,
			RouteID:     routes[2].ID,
			Status:      "on-route",
			Latitude:    19.0760,
			Longitude:   72.8777,
			LastUpdated: time.Now(),
		},

		// Bangalore Buses
		{
			BusNumber:   "KA-01-AB-5678",
			DriverID:    drivers[3].ID,
			RouteID:     routes[3].ID,
			Status:      "on-route",
			Latitude:    12.9716,
			Longitude:   77.5946,
			LastUpdated: time.Now(),
		},

		// Chennai Buses
		{
			BusNumber:   "TN-09-AB-9012",
			DriverID:    drivers[4].ID,
			RouteID:     routes[4].ID,
			Status:      "on-route",
			Latitude:    13.0629,
			Longitude:   80.2239,
			LastUpdated: time.Now(),
		},
	}

	for i := range buses {
		if err := s.DB.Create(&buses[i]).Error; err != nil {
			return err
		}
	}

	log.Println("Indian bus tracking seed data created successfully")
	return nil
}

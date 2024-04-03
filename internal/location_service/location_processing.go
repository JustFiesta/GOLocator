package main

import (
	"database/sql"
	"fmt"
	"golocator/models"
	"math"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "root:@tcp(localhost:3306)/golocator?charset=utf8mb4&parseTime=True&loc=Local"

	sqlDB, err := sql.Open("mysql", dsn)

	if err != nil {
		panic("failed to open mysql database")
	}

	db, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic("gorm failed to connect database")
	}
}

// getTraveledDistance - fetches and calculates all distance traveled since specified date
func getTraveledDistance(c echo.Context) error {
	// read http request parameters
	datetimeStr := c.QueryParam("datetime")

	fmt.Println("DEBUG datetimestr: ", datetimeStr)

	// validate date format
	if !validateDate(datetimeStr) {
		return c.JSON(http.StatusBadRequest, "Invalid date format")
	}

	// parse datetime
	datetime, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid date")
	}

	// fetch locations since the specified date
	var locations []models.Location
	if err := db.Where("date_time >= ?", datetime).Find(&locations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch locations")
	}

	// calculate total distance traveled
	totalDistance := 0.0
	for i := 0; i < len(locations)-1; i++ {
		totalDistance += haversine(locations[i].Latitude, locations[i].Longitude, locations[i+1].Latitude, locations[i+1].Longitude)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Total distance traveled since %s: %.2f kilometers", datetimeStr, totalDistance))
}

func validateDate(date string) bool {
	// ISO 8601 format
	iso8601Pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\+\d{2}:\d{2}|Z)$`

	// matching format with given string date
	match, err := regexp.MatchString(iso8601Pattern, date)

	if err != nil {
		return false
	}

	if !match {
		return false
	}

	// check if date is valid
	_, err = time.Parse(time.RFC3339, date)
	if err != nil {
		return false
	}

	return true
}

// haversine calculates the great-circle distance between two points on the Earth's surface
// given their longitudes and latitudes in degrees
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// haversine formula
	dlon := lon2Rad - lon1Rad
	dlat := lat2Rad - lat1Rad
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	r := 6371.0 // Earth radius in kilometers
	distance := r * c

	return distance
}

func main() {
	// initialize new DB connection with GORM
	initDB()

	// iniitlize new echo router
	e := echo.New()

	// logger from echo doc
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			fmt.Printf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

	// listen for GET requests on route /locationhistory/distancetraveled
	e.GET("/locationhistory/distancetraveled", getTraveledDistance)
	e.Logger.Fatal(e.Start(":1323"))
}

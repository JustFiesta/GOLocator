package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golocator/models"

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

// updateUserLocation - updates user location to specified coordinates
func updateUserLocation(c echo.Context) error {
	// get userid
	userID := c.Param("id")
	var user models.User

	// check if user is present in db
	if err := db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// create struct to hold latitude and longitude
	type Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	// parse request body into Coordinates struct
	var coordinates Coordinates
	if err := c.Bind(&coordinates); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	// validate coordinates
	latitude, longitude, err := validateCoordinates(fmt.Sprintf("%f,%f", coordinates.Latitude, coordinates.Longitude))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// create new user location
	location := models.Location{
		UserID:    user.ID,
		Latitude:  latitude,
		Longitude: longitude,
	}
	fmt.Println("DEBUG latitude: ", latitude, ", longitude: ", longitude)

	// persist location to db
	if err := db.Create(&location).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user location")
	}

	return c.JSON(http.StatusOK, "User location updated successfully")
}

// validateCoordinates - Validate user inputed coordinates AND return both latitude and longitude from one string
func validateCoordinates(coordStr string) (latitude, longitude float64, err error) {
	// Regexp pattern to validate latitude and longitude
	coordPattern := `^[-+]?([1-8]?\d(\.\d{1,8})?|90(\.0{1,8})?),\s*[-+]?(180(\.0{1,8})?|((1[0-7]\d)|([1-9]?\d))(\.\d{1,8})?)$`
	coordRegex := regexp.MustCompile(coordPattern)

	if !coordRegex.MatchString(coordStr) {
		return 0, 0, fmt.Errorf("Invalid coordinates format. Use decimal format (latitude,longitude)")
	}

	// Split coordinates string by comma
	parts := strings.Split(coordStr, ",")

	// Convert latitude and longitude to float64
	latitude, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid latitude")
	}

	longitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid longitude")
	}

	// Validate latitude range (-90 to 90) and longitude range (-180 to 180)
	if latitude < -90 || latitude > 90 {
		return 0, 0, fmt.Errorf("Latitude out of range (-90 to 90)")
	}

	if longitude < -180 || longitude > 180 {
		return 0, 0, fmt.Errorf("Longitude out of range (-180 to 180)")
	}

	return latitude, longitude, nil
}

// getUsersInLocation - fetches all users in given location
func getUsersInLocation(c echo.Context) error {
	// read http request parameters
	latStr := c.QueryParam("latitude")
	lonStr := c.QueryParam("longitude")
	radiusStr := c.QueryParam("radius")

	// convert params to corresponding types
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid latitude")
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid longitude")
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid radius")
	}

	// find users in location within a given radius
	var users []models.User
	if err := db.Table("location").Preload("Locations").Find(&users, "latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?", lat-radius, lat+radius, lon-radius, lon+radius).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to get users in location")
	}

	return c.JSON(http.StatusOK, users)
}

func main() {
	initDB()

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

	e.PUT("/location/:id", updateUserLocation)
	e.GET("/location/usersinlocation", getUsersInLocation)
	e.Logger.Fatal(e.Start(":1323"))
}

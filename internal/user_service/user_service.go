package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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
	fmt.Println("Debug: User id: ", userID)
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	fmt.Println("Debug: User struct: ", user)

	if err := db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user location")
	}

	return c.JSON(http.StatusOK, "User location updated successfully")
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

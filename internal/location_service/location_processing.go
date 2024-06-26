package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func runLocationService() {
	e := echo.New()
	e.GET("/locationhistory/distancetraveled", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func main() {
	runLocationService()
}
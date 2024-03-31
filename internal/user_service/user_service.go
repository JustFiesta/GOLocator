package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func runUserService() {
	e := echo.New()
	e.GET("/location/usersinlocation", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.PUT("/location/{user}", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func main() {
	runUserService()
}
package main

import (
    "github.com/labstack/echo"
    mw "github.com/labstack/echo/middleware"
)

func main() {
    // Echo instance
    e := echo.New()

    // Middleware
    e.Use(mw.RequestID())
    e.Use(mw.Logger())
    e.Use(mw.Recover())

    e.GET("/", EdgeController)

    // Start server
    e.Logger.Fatal(e.Start(":9001"))
}

package main

import (
  "github.com/labstack/echo"
  mw "github.com/labstack/echo/middleware"
)

func main() {
  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(mw.Logger())
  e.Use(mw.Recover())
  // e.Use(EdgeMiddleware)

  e.GET("/", EdgeController)

  // Start server
  e.Logger.Fatal(e.Start(":9001"))
}

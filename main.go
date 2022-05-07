package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zipzap11/pharm-API/controller"
	"github.com/zipzap11/pharm-API/db"
)

func main() {
	e := echo.New()
	db.InitDB()
	productController := controller.ProductController{}
	e.Use(middleware.Logger())
	e.GET("/products", productController.GetAllProducts)
	e.Start("localhost:8000")
}

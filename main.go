package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zipzap11/pharm-API/controller"
	"github.com/zipzap11/pharm-API/db"
	"github.com/zipzap11/pharm-API/repository"
	"github.com/zipzap11/pharm-API/usecase"
)

func main() {
	e := echo.New()
	db.InitDB()

	// product
	productRepository := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepository)
	productController := controller.NewProductController(productUsecase)

	// middleware
	e.Use(middleware.Logger())

	// path
	e.GET("/products", productController.GetAllProducts)

	e.Start("localhost:8000")
}

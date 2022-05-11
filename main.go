package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/config"
	"github.com/zipzap11/pharm-API/controller"
	"github.com/zipzap11/pharm-API/db"
	customMiddleware "github.com/zipzap11/pharm-API/middleware"
	"github.com/zipzap11/pharm-API/repository"
	"github.com/zipzap11/pharm-API/usecase"
	"github.com/zipzap11/pharm-API/util"
)

func main() {
	e := echo.New()
	db.InitDB()

	validator := validator.New()
	tokenProvider, err := util.NewPasetoProvider(config.GetSymmetricKey())
	if err != nil {
		logrus.Fatal(err)
	}

	// session
	sessionRepository := repository.NewSessionRepository(db.DB)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepository, tokenProvider)
	sessionController := controller.NewSessionController(sessionUsecase)

	// product
	productRepository := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepository)
	productController := controller.NewProductController(productUsecase)

	// category
	categoryRepository := repository.NewCategoryRepository(db.DB)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	categoryController := controller.NewCategoryController(categoryUsecase)

	// user
	userRepository := repository.NewUserRepository(db.DB)
	userUsecase := usecase.NewUserUsecase(userRepository, validator, tokenProvider, sessionRepository)
	userController := controller.NewUserController(userUsecase)

	// cart-item
	cartItemRepository := repository.NewCartItemRepository(db.DB)

	// cart
	cartRepository := repository.NewCartRepository(db.DB, cartItemRepository)
	cartUsecase := usecase.NewCartUsecase(cartRepository)
	cartController := controller.NewCartController(cartUsecase)

	// middleware
	e.Use(middleware.Logger())
	auth := e.Group("", customMiddleware.AuthPaseto(tokenProvider))

	// path
	auth.GET("/products", productController.GetAllProducts)
	e.GET("/products/:id", productController.FindById)
	e.GET("/categories", categoryController.GetAllCategories)
	e.POST("/users", userController.CreateUser)
	e.POST("/auth/login", userController.Login)
	e.POST("/auth/refresh", sessionController.RefreshSession)
	auth.GET("/carts", cartController.FindCart)
	e.Start("localhost:8000")
}

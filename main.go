package main

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
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

	var (
		midtransSnapClient snap.Client
		midtransCoreClient coreapi.Client
	)
	midtransSnapClient.New(config.GetMidtransAPIKey(), midtrans.Sandbox)
	midtransCoreClient.New(config.GetMidtransAPIKey(), midtrans.Sandbox)

	// session
	sessionRepository := repository.NewSessionRepository(db.DB)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepository, tokenProvider)
	sessionController := controller.NewSessionController(sessionUsecase)

	// category
	categoryRepository := repository.NewCategoryRepository(db.DB)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	categoryController := controller.NewCategoryController(categoryUsecase)

	// product
	productRepository := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepository, categoryRepository)
	productController := controller.NewProductController(productUsecase)

	// cart-item
	cartItemRepository := repository.NewCartItemRepository(db.DB)

	// cart
	cartRepository := repository.NewCartRepository(db.DB, cartItemRepository)
	cartUsecase := usecase.NewCartUsecase(cartRepository, cartItemRepository)
	cartController := controller.NewCartController(cartUsecase)

	// user
	userRepository := repository.NewUserRepository(db.DB)
	userUsecase := usecase.NewUserUsecase(userRepository, validator, tokenProvider, sessionRepository, db.DB, cartRepository)
	userController := controller.NewUserController(userUsecase)

	// address
	addressRepository := repository.NewAddressRepository(
		config.GetROAPIKey(),
		config.GetProvinceAPIUrl(),
		config.GetStateAPIUrl(),
		db.DB)
	addressUsecase := usecase.NewAddressUsecase(addressRepository)
	addressController := controller.NewAddressController(addressUsecase)

	// shipping
	shippingRepository := repository.NewShippingRepository(
		config.GetROAPIKey(),
		config.GetShippingOrigin(),
		config.GetROPriceURL(),
		db.DB)
	shippingUsecase := usecase.NewShippingUsecase(shippingRepository, cartItemRepository, cartRepository, addressRepository)
	shippingController := controller.NewShippingController(shippingUsecase)

	// transaction item
	transactionItemRepository := repository.NewTransactionItemRepository(db.DB)

	// transaction
	transactionRepository := repository.NewTransactionRepository(db.DB, transactionItemRepository)
	transactionUsecase := usecase.NewTransactionUsecase(cartUsecase, shippingUsecase, &midtransSnapClient, &midtransCoreClient, transactionRepository, userRepository, addressUsecase, db.DB, transactionItemRepository, cartItemRepository)
	transactionController := controller.NewTransactionController(transactionUsecase)

	// middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method} uri=${uri} status=${status} latency=${latency_human} header=${header:Authorization}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlExposeHeaders, echo.HeaderAuthorization},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			time.Sleep(time.Second * 1)
			return next(c)
		}
	})
	auth := e.Group("", customMiddleware.AuthPaseto(tokenProvider))
	// path
	e.GET("/products", productController.GetAllProducts)
	e.GET("/products/:id", productController.FindById)
	auth.POST("/products", productController.Create)
	auth.PUT("/products", productController.UpdateStock)
	auth.DELETE("/products/:id", productController.DeleteProduct)
	e.GET("/categories", categoryController.GetAllCategories)
	auth.POST("/categories", categoryController.CreateCategory)
	auth.GET("/users", userController.FindUsers)
	e.POST("/users", userController.CreateUser)
	//e.GET("/users:id", userController.FindById)
	e.POST("/users/super", userController.CreateSuperUser)
	e.POST("/auth/login", userController.Login)
	e.POST("/auth/refresh", sessionController.RefreshSession)
	e.GET("/auth/validate", sessionController.CheckSession)
	auth.GET("/auth/current-user", userController.FindCurrentUser)
	auth.GET("/carts", cartController.FindCart)
	auth.POST("/carts", cartController.AddItemToCart)
	auth.PUT("/carts", cartController.UpdateItemQuantity)
	auth.DELETE("/carts", cartController.RemoveItemFromCart)
	auth.GET("/provinces", addressController.GetProvinces)
	auth.GET("/states", addressController.GetStates)
	auth.POST("/addresses", addressController.CreateAddress)
	auth.GET("/addresses", addressController.GetAddressesByUserID)
	auth.GET("/shippings", shippingController.GetShippingPackages)
	auth.GET("/prices", transactionController.GetTotalPrice)
	auth.POST("/transactions", transactionController.CreateTransaction)
	e.POST("/transactions/callback", transactionController.HandleTransactionCallback)
	auth.GET("/transactions", transactionController.GetTransactionByUserID)
	auth.GET("/transactions/:id", transactionController.GetTransactionByID)

	e.Start("localhost:8000")
}

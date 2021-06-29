package routes

import (
	"net/http"
	"project/controllers"
	"project/middlewares"
	"project/models"

	"github.com/labstack/echo"
)

// Initialize echo router
func New() *echo.Echo {
	e := echo.New()
	// Root
	e.GET("/", rootHandler)

	// Auth routes
	authGroup := e.Group("")
	// Create new user & admin
	authGroup.POST("/register", controllers.CreateUsersController)
	// Login the user and get the access token
	authGroup.POST("/login", controllers.LoginUsersController)
	// Create refresh token handler (BONUS)
	authGroup.POST("/refresh", controllers.RefreshTokenController, middlewares.AuthenticationMiddleware(), middlewares.AuthorizationMiddleware("refresh"))

	// Admin Group
	adminGroup := e.Group("")
	// Use authentication and admin-level authorization
	adminGroup.Use(middlewares.AuthenticationMiddleware(), middlewares.AuthorizationMiddleware("basicAdminPermissions"))
	adminGroup.GET("/users", controllers.GetAllUsersController)
	adminGroup.GET("/user/:id", controllers.GetUserController)
	// Customers routes
	adminGroup.POST("/customer", controllers.CreateCustomerController)
	adminGroup.GET("/customers", controllers.GetAllCustomersController)
	adminGroup.GET("/customer/:id", controllers.GetCustomerController)
	adminGroup.PUT("/customer/:id", controllers.UpdateCustomerController)
	adminGroup.DELETE("/customer/:id", controllers.DeleteCustomerController)
	// Car routes
	adminGroup.POST("/car", controllers.CreateCarController)
	adminGroup.GET("/cars", controllers.GetAllCarsController)
	adminGroup.GET("/car/:id", controllers.GetCarController)
	adminGroup.PUT("/car/:id", controllers.UpdateCarController)
	adminGroup.DELETE("/car/:id", controllers.DeleteCarController)
	// Garage routes
	adminGroup.POST("/garage", controllers.CreateGarageController)
	adminGroup.GET("/garages", controllers.GetAllGaragesController)
	adminGroup.GET("/garage/:id", controllers.GetGarageController)
	adminGroup.PUT("/garage/:id", controllers.UpdateGarageController)
	adminGroup.DELETE("/garage/:id", controllers.DeleteGarageController)
	// Order routes
	adminGroup.POST("/order", controllers.CreateOrderController)
	adminGroup.GET("/orders", controllers.GetAllOrdersController)
	adminGroup.GET("/order/:id", controllers.GetOrderController)
	adminGroup.POST("/order/done", controllers.CloseOrderController)
	// Get all orders data based on the struct/database structure (BONUS)
	e.GET("/orders/v2", controllers.GetAllOrdersControllerV2)

	// Users Group (employees) only have acces to order api and read cars.
	userGroup := e.Group("")
	// Use authentification and user-level authorization
	userGroup.Use(middlewares.AuthenticationMiddleware(), middlewares.AuthorizationMiddleware("refresh"))
	userGroup.GET("/cars", controllers.GetAllCarsController)
	userGroup.GET("/car/:id", controllers.GetCarController)
	// Order routes
	userGroup.POST("/order", controllers.CreateOrderController)
	userGroup.GET("/orders", controllers.GetAllOrdersController)
	userGroup.GET("/order/:id", controllers.GetOrderController)
	userGroup.POST("/order/done", controllers.CloseOrderController)

	return e
}

// A functional easter-egg!
func rootHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Response{Message: "IT works!"})
}

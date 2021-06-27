package routes

import (
	"project/controllers"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// Customers routes
	e.POST("/customer", controllers.CreateCustomerController)
	// Pagination has been implemented
	// The default per_page size is 10
	//
	// e.g. We can also use query  -->  /customers?per_page=5&page=2
	//
	// It will return the items from index 6 to 10 (index >= 10)
	e.GET("/customers", controllers.GetAllCustomersController)
	// Filter by id
	e.GET("/customer/:id", controllers.GetCustomerController)
	// Update & delete by id
	e.PUT("/customer/:id", controllers.UpdateCustomerController)
	e.DELETE("/customer/:id", controllers.DeleteCustomerController)

	// Car routes
	e.POST("/car", controllers.CreateCarController)
	// Pagination has been implemented
	// The default per_page size is 10
	e.GET("/cars", controllers.GetAllCarsController)
	// Filter by id
	e.GET("/car/:id", controllers.GetCarController)
	// Update & delete by id
	e.PUT("/car/:id", controllers.UpdateCarController)
	e.DELETE("/car/:id", controllers.DeleteCarController)

	// Garage routes
	e.POST("/garage", controllers.CreateGarageController)
	// Pagination has been implemented
	// The default per_page size is 10
	e.GET("/garages", controllers.GetAllGaragesController)
	// Filter by id
	e.GET("/garage/:id", controllers.GetGarageController)
	// Update & delete by id
	e.PUT("/garage/:id", controllers.UpdateGarageController)
	e.DELETE("/garage/:id", controllers.DeleteGarageController)

	// Order routes
	// Request order
	e.POST("/order", controllers.CreateOrderController)
	// Get all the orders data optionally filtered by query params
	e.GET("/orders", controllers.GetAllOrdersController)
	// Filter by id
	e.GET("/order/:id", controllers.GetOrderController)
	// Close order
	e.POST("/order/done", controllers.CloseOrderController)

	// Get all orders data based on the struct/database structure (BONUS)
	e.GET("/orders/v2", controllers.GetAllOrdersControllerV2)
	return e
}

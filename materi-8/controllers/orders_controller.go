package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"project/config"
	"project/models"

	"github.com/labstack/echo"
)

/* POST /order --> to add order data
{
    "customer_id": 1,
    "car_id": 1,
    "with_driver": true,
    "estimated_days": 3
}
*/
func CreateOrderController(c echo.Context) error {
	order := models.Orders{}
	car := models.Cars{}
	c.Bind(&order)

	// Check how many inactive car available
	inactive_cars, err := car.CountActiveCars(config.DB, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Check whether selected car status is active
	car_active, err := car.IsACarActive(config.DB, order.Car_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Return resource error if all available cars are on active status
	if inactive_cars == 0 {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Car stock is empty",
			Error:   "Resource error",
		})
	} else if car_active {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Selected car is not available (booked / active)",
			Error:   "Resource error",
		})
	}
	// Store the request order data to the database
	added_order, err := order.SaveOrder(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Book or update the selected car to be non active on the database
	err = car.UpdateActiveCar(config.DB, added_order.Car_id, true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Calculate estimated price of an order
	estimated_price, err := added_order.CalculateEstPrice(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Preserve the order of data in JSON and return as specified on the requirement
	return c.JSON(http.StatusCreated, models.ComposedOrders{
		Order_id:        added_order.AddPrefixOrderId(added_order.ID), // add prefix
		Full_name:       *added_order.Customers.Full_name,
		Id_card:         *added_order.Customers.ID_card,
		Label:           *added_order.Cars.Label,
		Car_type:        *added_order.Cars.Car_type,
		Estimated_price: estimated_price,
	})
}

// GET /orders --> to get all the orders data and optionally filtered by query params
func GetAllOrdersController(c echo.Context) error {
	order := models.Orders{}
	// Check whether order data available
	counted_orders, err := order.CountAllOrders(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	if counted_orders == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no order data yet in database",
			Error:   err.Error(),
		})
	}
	// Load all order data as specified in the requirement
	meta, composed_orders, err := order.GetComposedAllOrders(config.DB, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.OrdersResponse{
		Meta: *meta,
		Data: *composed_orders,
	})
}

// GET /orders/v2 --> to get all orders data based on struct/database structure (BONUS)
func GetAllOrdersControllerV2(c echo.Context) error {
	order := models.Orders{}
	// Check whether order data available
	counted_orders, err := order.CountAllOrders(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	if counted_orders == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no order data yet in database",
			Error:   err.Error(),
		})
	}
	// Find all the available orders from database
	orders, err := order.GetAllOrders(config.DB, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get all the orders data successful",
		Data:    orders,
	})
}

// GET /order/:id --> to get an order data specified by id
func GetOrderController(c echo.Context) error {
	// Check the id parameter
	order_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Retreive Order object with the primary key (id)
	order := models.Orders{}
	the_order, err := order.GetOrder(config.DB, order_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: fmt.Sprintf("Get %s data succesful", order.AddPrefixOrderId(order_id)),
		Data:    the_order,
	})
}

/* POST /order --> close the status of an order to be done
{
    "order_id": "order-1"
}
*/
func CloseOrderController(c echo.Context) error {
	order := models.Orders{}
	prefix := models.Prefix{}
	c.Bind(&prefix)
	// Immediate response for a blank order_id
	prefix_id := prefix.Order_id
	if prefix_id == "" {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: `Order id is not found, (e.g. "order-1")`,
		})
	}
	// Trim prefix of the order id
	order_id, err := order.TrimPrefixOrderId(prefix_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: `order_id is invalid, use this (e.g. "order-1") instead`,
			Error:   err.Error(),
		})
	}
	// Check whether order is on the database
	closed_order, err := order.GetOrder(config.DB, order_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Order id is not found",
			Error:   err.Error(),
		})
	}
	// Check whether order id is done already
	done, err := closed_order.CheckOrderStatus(config.DB, order_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	if done != "on going" {
		return c.JSON(http.StatusFound, models.Response{
			Message: prefix.Order_id + " was done already",
		})
	}
	// Update total_days
	err = closed_order.UpdateTotalDays(config.DB, order_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Update the active status of the car of the order to be closed on the database
	err = closed_order.Cars.UpdateActiveCar(config.DB, closed_order.Cars.ID, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Update the late status of the order
	err = closed_order.UpdateLateStatus(config.DB, order_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Preserve the order of data in JSON
	return c.JSON(http.StatusOK, models.Response{
		Message: fmt.Sprintf("%s is marked done", order.AddPrefixOrderId(order_id)),
		Data: models.ComposedOrders{
			Order_id:    prefix.Order_id,
			Full_name:   *closed_order.Customers.Full_name,
			Id_card:     *closed_order.Customers.ID_card,
			Label:       *closed_order.Cars.Label,
			Car_type:    *closed_order.Cars.Car_type,
			Total_price: closed_order.TotalPrice(),
		},
	})
}

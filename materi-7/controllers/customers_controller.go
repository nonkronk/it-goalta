package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project/config"
	"project/models"

	"github.com/labstack/echo"
)

/* POST /customer --> to add customer data
{
    "full_name": "Irvan Tristian",
    "mobile": "0811223456",
    "address": "Mountain View",
    "email": "irvan.t@google.com",
    "id_card": "1708194517081945"
}
*/
func CreateCustomerController(c echo.Context) error {
	customer := models.Customers{}
	c.Bind(&customer)
	// Store the request data to the database
	added_customer, err := customer.SaveCustomer(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Customer data succesfully created",
		Data:    added_customer,
	})
}

// GET /customers --> to get all customers data
func GetAllCustomersController(c echo.Context) error {
	customer := models.Customers{}
	// Check whether customer data available
	counted_customers, err := customer.CountAllCustomers(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	if counted_customers == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no customer data yet in database",
			Error:   err.Error(),
		})
	}
	// Find all the available customers from database
	customers, err := customer.GetAllCustomers(config.DB, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get all the customers data successful",
		Data:    customers,
	})
}

// GET /customer/:id --> to get a customer data specified by id
func GetCustomerController(c echo.Context) error {
	// Check the id parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Retreive customer object with the primary key (id)
	customer := models.Customers{}
	the_customer, err := customer.GetCustomer(config.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a customer data succesful",
		Data:    the_customer,
	})
}

// PUT /customer/:id --> to update a customer data specified by id
func UpdateCustomerController(c echo.Context) error {
	// Create a map to support attributes update
	// and populate json data.
	// It will only update non-zero value fields
	request_body := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&request_body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	// Check id parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Update customer data in database
	customer := models.Customers{}
	updated_customer, err := customer.UpdateCustomer(config.DB, id, request_body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Update a customer data succesful",
		Data:    updated_customer,
	})
}

// DELETE /customer/:id --> to delete a customer data specified by id from database
func DeleteCustomerController(c echo.Context) error {
	// Check id parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Delete requested data off the database
	customer := models.Customers{}
	deleted_customer, err := customer.DeleteCustomer(config.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Delete a customer data succesful",
		Data:    deleted_customer,
	})
}

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
	if err := config.DB.Debug().Create(&customer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Customer data succesfully created",
		Data:    customer,
	})
}

// GET /customers --> to get all customers data
func GetAllCustomersController(c echo.Context) error {
	customer := models.Customers{}
	customers := []models.Customers{}
	// Check whether customer data available
	counted_customers, err := customer.CountAllCustomers(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if counted_customers == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no customer data yet in database",
			Status:  err.Error(),
		})
	}
	// With pagination implemented
	// Find all the available customers from database
	if err := config.DB.Debug().Scopes(models.Paginate(c)).Find(&customers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
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
	customer_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Retreive customer object with the primary key (id)
	customer := models.Customers{}
	if err := config.DB.Debug().First(&customer, customer_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a customer data succesful",
		Data:    customer,
	})
}

// PUT /customer/:id --> to update a customer data specified by id
func UpdateCustomerController(c echo.Context) error {
	// Create a map to support attributes update
	// and populate json data.
	// It will only update non-zero value fields
	customer := models.Customers{}
	request_body := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&request_body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	// Check id parameter
	customer_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Update customer data in database
	if err := config.DB.Debug().Model(&models.Customers{}).Where("id = ?", customer_id).Updates(request_body).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Get the selected customer object to prevent null or missing data
	// so that the method return complete customer data
	if err := config.DB.Debug().First(&customer, customer_id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Update a customer data succesful",
		Data:    customer,
	})
}

// DELETE /customer/:id --> to delete a customer data specified by id from database
func DeleteCustomerController(c echo.Context) error {
	// Check id parameter
	customer_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Delete requested data off the database
	var customer models.Customers
	var deleted_customer models.Customers
	if err := config.DB.Debug().First(&deleted_customer, customer_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if err := config.DB.Debug().Delete(&customer, customer_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Delete a customer data succesful",
		Data:    deleted_customer})
}

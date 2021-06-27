package controllers

import (
	"net/http"
	"strconv"

	"project/config"
	"project/models"

	"github.com/labstack/echo"
)

/* POST /car --> to add car data
{
    "garage_id": 1,
    "car_type": "City Car",
    "car": "Suzuki Swift",
    "transmission": "MT",
    "fuel": "Gasoline",
	"price_per_day": 300000
}
*/
func CreateCarController(c echo.Context) error {
	car := models.Cars{}
	c.Bind(&car)
	// Store the request data to the database
	added_car, err := car.SaveCar(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Car data is invalid or garage_id is not recognized",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Car data succesfully created",
		Data:    added_car,
	})
}

// GET /cars --> to get all cars data
func GetAllCarsController(c echo.Context) error {
	car := models.Cars{}
	// Check whether car data available
	counted_cars, err := car.CountAllCars(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	if counted_cars == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no car data yet in database",
			Error:   err.Error(),
		})
	}
	// Find all the available cars from database
	cars, err := car.GetAllCars(config.DB, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get all the cars data successful",
		Data:    cars,
	})
}

// GET /car/:id --> to get a car data specified by id
func GetCarController(c echo.Context) error {
	// Check the id parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Retreive car object with the primary key (id)
	car := models.Cars{}
	the_car, err := car.GetCar(config.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a car data succesful",
		Data:    the_car,
	})
}

// PUT /car/:id --> to update a car data specified by id
func UpdateCarController(c echo.Context) error {
	// Fix the inconsistent requirement on the key name of the car
	// label --> car
	car := models.Cars{}
	request_body, err := car.FixInconsistentKey(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
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
	// Update garage data in database
	updated_car, err := car.UpdateCar(config.DB, id, request_body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Update a garage data succesful",
		Data:    updated_car,
	})
}

// DELETE /car/:id --> to delete a car data specified by id from database
func DeleteCarController(c echo.Context) error {
	// Check id parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Delete requested data off the database
	car := models.Cars{}
	deleted_car, err := car.DeleteCar(config.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "We can't delete a car that has been ordered in the database",
			Error:   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Delete a car data succesful",
		Data:    deleted_car})
}

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
	if err := config.DB.Debug().Create(&car).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	// Populate garage data; the foreignkey of a car
	if err := config.DB.Debug().Model(&models.Garages{}).Where("id = ?", car.Garage_id).Take(&car.Garages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Car data succesfully created",
		Data:    car,
	})
}

// GET /cars --> to get all cars data
func GetAllCarsController(c echo.Context) error {
	car := models.Cars{}
	cars := []models.Cars{}
	// Check whether car data available
	counted_cars, err := car.CountAllCars(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if counted_cars == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no car data yet in database",
			Status:  err.Error(),
		})
	}
	// With pagination implemented
	if err := config.DB.Debug().Scopes(models.Paginate(c)).Find(&cars).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	for i := range cars {
		if err := config.DB.Debug().Model(&models.Garages{}).Where("id = ?", cars[i].Garage_id).Take(&cars[i].Garages).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get all the cars data successful",
		Data:    cars,
	})
}

// GET /car/:id --> to get a car data specified by id
func GetCarController(c echo.Context) error {
	// Check the id parameter
	car_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Retreive car object with the primary key (id)
	car := models.Cars{}
	if err := config.DB.Debug().First(&car, car_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if err := config.DB.Debug().Model(&models.Garages{}).Where("id = ?", car.Garage_id).Take(&car.Garages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a car data succesful",
		Data:    car,
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
			Status: err.Error(),
		})
	}
	// Check id parameter
	car_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Update car data in database
	if err := config.DB.Debug().Model(&models.Cars{}).Where("id = ?", car_id).Updates(request_body).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Get the selected car object to prevent null or missing data
	// so that the method return complete car data
	if err := config.DB.Debug().First(&car, car_id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	if err := config.DB.Debug().Model(&models.Garages{}).Where("id = ?", car.Garage_id).Take(&car.Garages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Update a car data succesful",
		Data:    car,
	})
}

// DELETE /car/:id --> to delete a car data specified by id from database
func DeleteCarController(c echo.Context) error {
	// Check id parameter
	car_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Delete requested data off the database
	var car models.Cars
	var deleted_car models.Cars
	if err := config.DB.Debug().First(&deleted_car, car_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if err := config.DB.Debug().Model(&models.Garages{}).Where("id = ?", deleted_car.Garage_id).Take(&deleted_car.Garages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	config.DB.Debug().Model(&models.Orders{}).Association("Cars").Delete(models.Cars{})
	if err := config.DB.Debug().Delete(&car, car_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "We can't delete a car that has been ordered in the database",
			Status:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Delete a car data succesful",
		Data:    deleted_car})
}

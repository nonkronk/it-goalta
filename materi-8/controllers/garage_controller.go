package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project/config"
	"project/models"

	"github.com/labstack/echo"
)

/* POST /garage --> to add garage data
{
    "owner": "Irvan Tristian",
    "address": "Mountain View",
    "mobile": "0811223456"
}
*/
func CreateGarageController(c echo.Context) error {
	garage := models.Garages{}
	c.Bind(&garage)
	// Store the request data to the database
	if err := config.DB.Debug().Create(&garage).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Garage data succesfully created",
		Data:    garage,
	})
}

// GET /garages --> to get all garages data
func GetAllGaragesController(c echo.Context) error {
	garage := models.Garages{}
	garages := []models.Garages{}
	// Check whether garage data available
	counted_garages, err := garage.CountAllGarages(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if counted_garages == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no garage data yet in database",
			Status:  err.Error(),
		})
	}
	// With pagination implemented
	// Find all the available garages from database
	if err := config.DB.Debug().Scopes(models.Paginate(c)).Find(&garages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get all the garages data successful",
		Data:    garages,
	})
}

// GET /garage/:id --> to get a garage data specified by id
func GetGarageController(c echo.Context) error {
	// Check the id parameter
	garage_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Retreive garage object with the primary key (id)
	garage := models.Garages{}
	if err := config.DB.Debug().First(&garage, garage_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a garage data succesful",
		Data:    garage,
	})
}

// PUT /garage/:id --> to update a garage data specified by id
func UpdateGarageController(c echo.Context) error {
	// Create a map to support attributes update
	// and populate json data.
	// It will only update non-zero value fields
	garage := models.Garages{}
	request_body := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&request_body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	// Check id parameter
	garage_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Update garage data in database
	if err := config.DB.Debug().Model(&models.Garages{}).Where("id = ?", garage_id).Updates(request_body).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Get the selected garage object to prevent null or missing data
	// so that the method return complete garage data
	if err := config.DB.Debug().First(&garage, garage_id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Update a garage data succesful",
		Data:    garage,
	})
}

// DELETE /garage/:id --> to delete a garage data specified by id from database
func DeleteGarageController(c echo.Context) error {
	// Check id parameter
	garage_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Delete requested data off the database
	var garage models.Garages
	var deleted_garage models.Garages
	if err := config.DB.Debug().First(&deleted_garage, garage_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if err := config.DB.Debug().Delete(&garage, garage_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Delete a garage data succesful",
		Data:    deleted_garage})
}

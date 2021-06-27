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
	added_garage, err := garage.SaveGarage(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Garage data succesfully created",
		Data:    added_garage,
	})
}

// GET /garages --> to get all garages data
func GetAllGaragesController(c echo.Context) error {
	garage := models.Garages{}
	// Check whether garage data available
	counted_garages, err := garage.CountAllGarages(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	if counted_garages == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no garage data yet in database",
			Error:   err.Error(),
		})
	}
	// Find all the available garages from database
	garages, err := garage.GetAllGarages(config.DB, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Retreive garage object with the primary key (id)
	garage := models.Garages{}
	the_garage, err := garage.GetAGarage(config.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a garage data succesful",
		Data:    the_garage,
	})
}

// PUT /garage/:id --> to update a garage data specified by id
func UpdateGarageController(c echo.Context) error {
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
	// Update garage data in database
	garage := models.Garages{}
	updated_garage, err := garage.UpdateGarage(config.DB, id, request_body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Update a garage data succesful",
		Data:    updated_garage,
	})
}

// DELETE /garage/:id --> to delete a garage data specified by id from database
func DeleteGarageController(c echo.Context) error {
	// Check id parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
	}
	// Delete requested data off the database
	garage := models.Garages{}
	deleted_garage, err := garage.DeleteGarage(config.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Delete a garage data succesful",
		Data:    deleted_garage,
	})
}

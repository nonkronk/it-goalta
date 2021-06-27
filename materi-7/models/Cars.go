package models

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

/* Add a Car POST request example:
{
    "garage_id": 1,
    "car_type": "City Car",
    "car": "Suzuki Swift",
    "transmission": "MT",
    "fuel": "Gasoline",
	"price_per_day": 300000
}
*/
type Cars struct {
	ID            int       `json:"car_id" form:"car_id" gorm:"primaryKey;autoIncrement;not null"`
	Garage_id     int       `json:"garage_id" form:"garage_id" gorm:"not null"`
	Car_type      *string   `json:"car_type" form:"car_type" gorm:"not null"`
	Label         *string   `json:"car" form:"car" gorm:"not null"`
	Transmission  *string   `json:"transmission" form:"transmission" gorm:"not null"`
	Fuel          *string   `json:"fuel" form:"fuel" gorm:"not null"`
	Price_per_day int       `json:"price_per_day" form:"price_per_day" gorm:"not null"`
	Is_active     bool      `json:"is_active" form:"is_active" gorm:"default:true"`
	Created_at    time.Time `json:"created_at" form:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
	// Use Garage_id as a foreignKey
	Garages *Garages `json:",omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:Garage_id"`
}

// Store the request data to database
func (c *Cars) SaveCar(db *gorm.DB) (*Cars, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Cars{}, err
	}
	err = c.PopulateCarGarage(db, c.Garage_id)
	return c, err
}

// Populate garage data; the foreignkey of a car
func (c *Cars) PopulateCarGarage(db *gorm.DB, garage_id int) error {
	err := db.Debug().Model(&Garages{}).Where("id = ?", garage_id).Take(&c.Garages).Error
	return err
}

// Find all the available cars from database
func (c *Cars) GetAllCars(db *gorm.DB, ec echo.Context) (*[]Cars, error) {
	cars := []Cars{}
	// With pagination implemented
	err := db.Debug().Scopes(Paginate(ec)).Find(&cars).Error
	if err != nil {
		return &[]Cars{}, err
	}
	for i := range cars {
		err = db.Debug().Model(&Garages{}).Where("id = ?", cars[i].Garage_id).Take(&cars[i].Garages).Error
		if err != nil {
			return &[]Cars{}, err
		}
	}
	return &cars, err
}

// Retreive car object with the primary key (id)
func (c *Cars) GetCar(db *gorm.DB, car_id int) (*Cars, error) {
	err := db.Debug().First(&c, car_id).Error
	if err != nil {
		return &Cars{}, err
	}
	err = db.Debug().Model(&Garages{}).Where("id = ?", c.Garage_id).Take(&c.Garages).Error
	return c, err
}

// Update car data in database
func (c *Cars) UpdateCar(db *gorm.DB, car_id int, request_body map[string]interface{}) (*Cars, error) {

	err := db.Debug().Model(&Cars{}).Where("id = ?", car_id).Updates(request_body).Error
	if err != nil {
		return &Cars{}, err
	}
	// Get the selected car object to prevent null or missing data
	// so that the method return complete car data
	err = db.Debug().First(&c, car_id).Error
	if err != nil {
		return &Cars{}, err
	}
	err = c.PopulateCarGarage(db, c.Garage_id)
	return c, err
}

// Delete car data from database
func (c *Cars) DeleteCar(db *gorm.DB, car_id int) (*Cars, error) {
	var deleted_car Cars
	err := db.Debug().First(&deleted_car, car_id).Error
	if err != nil {
		return &Cars{}, err
	}
	err = db.Debug().Model(&Garages{}).Where("id = ?", deleted_car.Garage_id).Take(&deleted_car.Garages).Error
	if err != nil {
		return &Cars{}, err
	}
	db.Debug().Model(&Orders{}).Association("Cars").Delete(Cars{})
	if err != nil {
		return &Cars{}, err
	}
	err = db.Debug().Delete(&c, car_id).Error
	return &deleted_car, err
}

// Fix the inconsistent requirement of the car struct/table
func (c *Cars) FixInconsistentKey(ec echo.Context) (map[string]interface{}, error) {
	// Create a map to support attributes update
	// and populate json data.
	// It will only update non-zero value fields
	request_body := make(map[string]interface{})
	err := json.NewDecoder(ec.Request().Body).Decode(&request_body)
	if err != nil {
		return make(map[string]interface{}), err
	}
	label := request_body["car"]
	delete(request_body, "car")
	request_body["label"] = label
	return request_body, err
}

// Count the number of car objects available in database
func (c *Cars) CountAllCars(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Debug().Model(&Cars{}).Count(&count).Error
	return count, err
}

// Count the number of active cars available in database
func (c *Cars) CountActiveCars(db *gorm.DB, status bool) (int64, error) {
	var count int64
	err := db.Debug().Model(&Cars{}).Where("is_active = ?", status).Count(&count).Error
	return count, err
}

// Check whether selected car status is active
func (c *Cars) IsACarActive(db *gorm.DB, car_id int) (bool, error) {
	var result bool
	err := db.Debug().Model(&Cars{}).Select("is_active").Where("id = ?", car_id).Scan(&result).Error
	return result, err
}

// Update the active status of the car of an order
func (c *Cars) UpdateActiveCar(db *gorm.DB, car_id int, status bool) error {
	err := db.Debug().Model(&Cars{}).Where("id = ?", car_id).Update("is_active", status).Error
	return err
}

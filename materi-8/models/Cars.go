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

// Fix the inconsistent requirement of the car struct/table
func (c *Cars) FixInconsistentKey(ec echo.Context) (map[string]interface{}, error) {
	// Create a map to support attributes update
	// and populate json data.
	// It will only update non-zero value fields
	request_body := make(map[string]interface{})
	if err := json.NewDecoder(ec.Request().Body).Decode(&request_body); err != nil {
		return make(map[string]interface{}), err
	}
	label := request_body["car"]
	delete(request_body, "car")
	request_body["label"] = label
	return request_body, nil
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

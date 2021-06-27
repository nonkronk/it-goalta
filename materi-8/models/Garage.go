package models

import (
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Garages struct {
	ID         int       `json:"garage_id" form:"garage_id" gorm:"primaryKey;autoIncrement;not null"`
	Owner      *string   `json:"owner" form:"owner" gorm:"not null"`
	Address    *string   `json:"address" form:"address" gorm:"not null"`
	Mobile     *string   `json:"mobile" form:"mobile" gorm:"not null"`
	Created_at time.Time `json:"created_at" form:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
}

// Store the request data to database
func (g *Garages) SaveGarage(db *gorm.DB) (*Garages, error) {
	err := db.Debug().Create(&g).Error
	return g, err
}

// Find all the available garages from database
func (g *Garages) GetAllGarages(db *gorm.DB, c echo.Context) (*[]Garages, error) {
	garages := []Garages{}
	// With pagination implemented
	err := db.Debug().Scopes(Paginate(c)).Find(&garages).Error
	return &garages, err
}

// Retreive garage object with the primary key (id)
func (g *Garages) GetAGarage(db *gorm.DB, garage_id int) (*Garages, error) {
	err := db.Debug().First(&g, garage_id).Error
	return g, err
}

// Update garage data in database
func (g *Garages) UpdateGarage(db *gorm.DB, garage_id int, request_body map[string]interface{}) (*Garages, error) {
	if err := db.Debug().Model(&Garages{}).Where("id = ?", garage_id).Updates(request_body).Error; err != nil {
		return &Garages{}, err
	}
	// Get the selected garage object to prevent null or missing data
	// so that the method return complete garage data
	err := db.Debug().First(&g, garage_id).Error
	return g, err
}

// Delete garage data from database
func (g *Garages) DeleteGarage(db *gorm.DB, garage_id int) (*Garages, error) {
	deleted_garage := Garages{}
	if err := db.Debug().First(&deleted_garage, garage_id).Error; err != nil {
		return &Garages{}, err
	}
	err := db.Debug().Delete(&g, garage_id).Error
	return &deleted_garage, err
}

// Count the number of garage objects available in database
func (g *Garages) CountAllGarages(db *gorm.DB) (int64, error) {
	result := db.Debug().First(&g)
	err := result.Error
	return result.RowsAffected, err
}

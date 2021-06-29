package models

import (
	"time"

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

// Count the number of garage objects available in database
func (g *Garages) CountAllGarages(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Debug().Model(&Garages{}).Count(&count).Error
	return count, err
}

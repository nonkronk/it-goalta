package models

import (
	"time"

	"gorm.io/gorm"
)

/* Add a Customer POST request example:
{
    "full_name": "Irvan Tristian",
    "mobile": "0811223456",
    "address": "Mountain View",
    "email": "irvan.t@google.com",
    "id_card": "1708194517081945"
}
*/
type Customers struct {
	ID         int       `json:"customer_id" form:"customer_id" gorm:"primaryKey;autoIncrement;not null"`
	Full_name  *string   `json:"full_name" form:"full_name" gorm:"not null"`
	Mobile     *string   `json:"mobile" form:"mobile" gorm:"not null"`
	Address    *string   `json:"address" form:"address"`
	Email      *string   `json:"email" form:"email"`
	ID_card    *string   `json:"id_card" query:"id_card" form:"id_card" gorm:"not null"`
	Created_at time.Time `json:"created_at" form:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
}

// Count the number of car objects available in database
func (c *Customers) CountAllCustomers(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Debug().Model(&Customers{}).Count(&count).Error
	return count, err
}

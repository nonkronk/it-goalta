package models

import (
	"time"

	"github.com/labstack/echo"
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

// Store the request data to database
func (c *Customers) SaveCustomer(db *gorm.DB) (*Customers, error) {
	err := db.Debug().Create(&c).Error
	return c, err
}

// Find all the available customers from database
func (c *Customers) GetAllCustomers(db *gorm.DB, ce echo.Context) (*[]Customers, error) {
	customers := []Customers{}
	// With pagination implemented
	err := db.Debug().Scopes(Paginate(ce)).Find(&customers).Error
	return &customers, err
}

// Retreive customer object with the primary key (id)
func (c *Customers) GetCustomer(db *gorm.DB, customer_id int) (*Customers, error) {
	err := db.Debug().First(&c, customer_id).Error
	return c, err
}

// Update customer data in database
func (c *Customers) UpdateCustomer(db *gorm.DB, customer_id int, request_body map[string]interface{}) (*Customers, error) {
	if err := db.Debug().Model(&Customers{}).Where("id = ?", customer_id).Updates(request_body).Error; err != nil {
		return &Customers{}, err
	}
	// Get the selected customer object to prevent null or missing data
	// so that the method return cosmplete customer data
	err := db.Debug().First(&c, customer_id).Error
	return c, err
}

// Count the number of car objects available in database
func (c *Customers) CountAllCustomers(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Debug().Model(&Customers{}).Count(&count).Error
	return count, err
}

// Delete customer data from database
func (c *Customers) DeleteCustomer(db *gorm.DB, customer_id int) (*Customers, error) {
	deleted_customer := Customers{}
	if err := db.Debug().First(&deleted_customer, customer_id).Error; err != nil {
		return &Customers{}, err
	}
	err := db.Debug().Delete(&c, customer_id).Error
	return &deleted_customer, err
}

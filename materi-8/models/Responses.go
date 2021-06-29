package models

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Data transfer object to preserve order of JSON key value when GET /order requested
type ComposedOrders struct {
	Id              int    `json:"-"`
	Order_id        string `json:"order_id,omitempty"`
	Full_name       string `json:"fullname,omitempty"`
	Id_card         string `json:"id_card,omitempty"`
	Mobile          string `json:"mobile,omitempty"`
	Label           string `json:"car,omitempty"`
	Car_type        string `json:"car_type,omitempty"`
	Total_days      int    `json:"days,omitempty"`
	Estimated_days  int    `json:"estimated_days,omitempty"`
	With_driver     bool   `json:"with_driver,omitempty"`
	Status          string `json:"is_late,omitempty"`
	Estimated_price int    `json:"estimated_price,omitempty"`
	Total_price     int    `json:"total_price,omitempty"`
}

type OrdersResponse struct {
	Meta Meta             `json:"meta"`
	Data []ComposedOrders `json:"data"`
}

type Meta struct {
	Total_data int64 `json:"total_data"`
	Per_page   int   `json:"per_page"`
	Page       int   `json:"page"`
}

type Response struct {
	Message               string      `json:"message,omitempty"`
	Data                  interface{} `json:"data,omitempty"`
	Status                interface{} `json:"status,omitempty"`
	AccessToken           string      `json:"access_token,omitempty"`
	AccessTokenExpiresAt  int64       `json:"access_token_expires_at,omitempty"`
	RefreshToken          string      `json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt int64       `json:"refresh_token_expires_at,omitempty"`
}

type Prefix struct {
	Order_id string `json:"order_id" form:"order_id" gorm:"not null"`
}

type FieldError struct {
	FieldName    string `json:"field"`
	ErrorMessage string `json:"message"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

// Calculate Metadata of the orders
func (m *Meta) GetResult(c echo.Context, result int64) {
	m.Total_data = result
	actualpage := (m.Total_data / 10) + 1
	// Add default page & per_page size (10)
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page == 0 || int64(page) > actualpage {
		page = 1
	}
	page_size, _ := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil || page_size == 0 {
		page_size = 10 // Default page_size params
	}
	m.Page = page
	m.Per_page = page_size
}

// Add prefix to the order id
func (o *ComposedOrders) AddPrefixOrderId(order_id int) string {
	str_id := strconv.Itoa(order_id)
	prefix_id := fmt.Sprintf("%s"+str_id, "order-") // add prefix
	return prefix_id
}

// Implement pagination
func Paginate(c echo.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}
		per_page, _ := strconv.Atoi(c.QueryParam("per_page"))
		switch {
		case per_page > 100:
			per_page = 100
		case per_page <= 0:
			per_page = 10
		}
		offset := (page - 1) * per_page
		return db.Offset(offset).Limit(per_page)
	}
}

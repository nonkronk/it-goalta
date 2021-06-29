package models

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Orders struct {
	ID             int  `json:"order_id" query:"order_id" form:"order_id" gorm:"primaryKey;autoIncrement;not null"`
	Customer_id    int  `json:"customer_id" form:"customer_id"`
	Car_id         int  `json:"car_id" form:"car_id"`
	Estimated_days int  `json:"estimated_days" form:"estimated_days" gorm:"not null"`
	With_driver    bool `json:"with_driver" form:"with_driver" gorm:"type:boolean"`

	// Need to store total_days in the database to maintain its integrity
	Total_days int `json:"total_days"`

	Status     string    `json:"is_late" gorm:"default:on going"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Include Custumer_id and Car_id as foreign key
	Customers *Customers `json:",omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:Customer_id"`
	Cars      *Cars      `json:",omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:Car_id"`
}

// Store the request data to database
func (o *Orders) SaveOrder(db *gorm.DB) (*Orders, error) {
	if err := o.PopulateOrderCustomer(db, o.Customer_id); err != nil {
		return &Orders{}, err
	}
	if err := o.PopulateOrderCar(db, o.Car_id); err != nil {
		return &Orders{}, err
	}
	err := db.Debug().Create(&o).Error
	return o, err
}

// Get order data by id
func (o *Orders) GetOrder(db *gorm.DB, order_id int) (*Orders, error) {
	if err := db.Debug().First(&o, order_id).Error; err != nil {
		return &Orders{}, err
	}
	if err := o.PopulateOrderCustomer(db, o.Customer_id); err != nil {
		return &Orders{}, err
	}
	err := o.PopulateOrderCar(db, o.Car_id)
	return o, err
}

// Find all the available cars from database
func (o *Orders) GetComposedAllOrders(db *gorm.DB, c echo.Context) (*Meta, *[]ComposedOrders, error) {
	composed_orders := []ComposedOrders{}
	meta := Meta{}
	order_id := c.QueryParam("order_id")
	is_late := c.QueryParam("is_late")
	id_card := c.QueryParam("id_card")

	// With pagination implemented, execute query dynamically from parameters
	chain := db.Debug().Model(&Orders{}).
		Select(
			`orders.id`,
			`customers.full_name`,
			`customers.id_card`,
			`customers.mobile`,
			`cars.label`,
			`cars.car_type`,
			`orders.total_days`,
			`orders.estimated_days`,
			`orders.with_driver`,
			`orders.status`).
		Joins("left join customers on customers.id = orders.customer_id").
		Joins("left join cars on cars.id = orders.car_id")
	if order_id != "" {
		chain = chain.Where("orders.id = " + order_id)
	}
	if is_late != "" {
		chain = chain.Where(fmt.Sprintf("orders.status = '%s'", is_late))
	}
	if id_card != "" {
		chain = chain.Where("customers.id_card = " + id_card)
	}
	result := chain.Scan(&composed_orders)
	if err := result.Error; err != nil {
		return &Meta{}, &[]ComposedOrders{}, err
	}
	// Get the metadata
	meta.GetResult(c, result.RowsAffected)
	// Paginate the result
	if err := chain.Scopes(Paginate(c)).Scan(&composed_orders).Error; err != nil {
		return &Meta{}, &[]ComposedOrders{}, err
	}
	// Manipulate and add prefix to each order ID
	for i := range composed_orders {
		str_id := strconv.Itoa(composed_orders[i].Id)
		prefix_id := fmt.Sprintf("%s"+str_id, "order-") // add prefix
		composed_orders[i].Order_id = prefix_id
	}
	return &meta, &composed_orders, nil
}

// Get all the orders data and optionally filtered by query params
func (o *Orders) GetAllOrders(db *gorm.DB, c echo.Context) (*[]Orders, error) {
	orders := []Orders{}
	// With pagination implemented
	if err := db.Debug().Scopes(Paginate(c)).Find(&orders).Error; err != nil {
		return &[]Orders{}, err
	}
	for i := range orders {
		if err := db.Debug().Model(&Customers{}).Where("id = ?", orders[i].Customer_id).Take(&orders[i].Customers).Error; err != nil {
			return &[]Orders{}, err
		}
		if err := db.Debug().Model(&Cars{}).Where("id = ?", orders[i].Car_id).Take(&orders[i].Cars).Error; err != nil {
			return &[]Orders{}, err
		}
	}
	if err := db.Debug().Scopes(Paginate(c)).Find(&orders).Error; err != nil {
		return &[]Orders{}, err
	}
	return &orders, nil
}

// Populate customer data; the foreignkey of an order
func (o *Orders) PopulateOrderCustomer(db *gorm.DB, customer_id int) error {
	err := db.Debug().Model(&Customers{}).Where("id = ?", customer_id).Take(&o.Customers).Error
	return err
}

// Populate car data; the foreignkey of an order
func (o *Orders) PopulateOrderCar(db *gorm.DB, car_id int) error {
	err := db.Debug().Model(&Cars{}).Where("id = ?", car_id).Take(&o.Cars).Error
	return err
}

// Calculate estimated price of an order
func (o *Orders) CalculateEstPrice(db *gorm.DB) (int, error) {
	if err := o.PopulateOrderCar(db, o.Car_id); err != nil {
		return 0, err
	}
	driver_price := 150000 // Hard coded driver price
	if !o.With_driver {
		driver_price = 0
	}
	p := float64(o.Cars.Price_per_day)
	e := float64(o.Estimated_days)
	d := float64(driver_price)
	return int((p * e) + (d * e)), nil
}

// Count the number of order objects available in database
func (o *Orders) CountAllOrders(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Debug().Model(&Orders{}).Count(&count).Error
	return count, err
}

// Add prefix to the order id
func (o *Orders) AddPrefixOrderId(order_id int) string {
	str_id := strconv.Itoa(order_id)
	prefix_id := fmt.Sprintf("%s"+str_id, "order-") // add prefix
	return prefix_id
}

// Trim prefix of the order id
func (o *Orders) TrimPrefixOrderId(prefix_id string) (int, error) {
	trimmed_id := bytes.TrimPrefix([]byte(prefix_id), []byte("order-"))
	order_id, err := strconv.Atoi(string(trimmed_id))
	return order_id, err
}

// Check the status of an order by id whether is it late (yes or no)
func (o *Orders) CheckOrderStatus(db *gorm.DB, order_id int) (string, error) {
	var is_late string
	err := db.Debug().Model(&Orders{}).Select("status").Where("id = ?", order_id).Scan(&is_late).Error
	return is_late, err
}

// Update the late status of all the orders
func (o *Orders) UpdateLateStatus(db *gorm.DB, order_id int) error {
	if o.Total_days > o.Estimated_days {
		err := db.Debug().Model(&Orders{}).Where("id = ?", order_id).Update("status", "yes").Error
		return err
	}
	err := db.Debug().Model(&Orders{}).Where("id = ?", order_id).Update("status", "no").Error
	return err
}

// Determine total_days of an order occured
func (o *Orders) UpdateTotalDays(db *gorm.DB, order_id int) error {
	var created_at time.Time
	if err := db.Debug().Model(&Orders{}).Select("Created_at").Where("id = ?", order_id).Scan(&created_at).Error; err != nil {
		return err
	}
	total_days := time.Now().Day() - created_at.Day()
	if total_days < 1 {
		total_days = 1
	}
	err := db.Debug().Model(&Orders{}).Where("id = ?", order_id).Update("total_days", total_days).Error
	return err
}

// Calculate total price of an order
func (o *Orders) TotalPrice() int {
	price_per_day := o.Cars.Price_per_day
	driver_price := 150000 // Hard coded driver price
	if !o.With_driver {
		driver_price = 0
	}
	p := float64(price_per_day)
	t := float64(o.Total_days)
	d := float64(driver_price)
	return int((p * t) + (d * t))
}

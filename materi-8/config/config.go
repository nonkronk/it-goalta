package config

import (
	"fmt"
	"log"
	"os"
	"project/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// It loads values from .env into the system
	// Is it the best practice?
	// I include mysql configuration on the .env file. Check the .env out!
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
}

func InitDB() {
	// Connect to the database using DB configuration from .env
	// The .env file must be adjusted to comply with your database configuration
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	connectDB(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName)
	initMigrate()
}

func initMigrate() {
	// Automigrate schema from the struct models
	DB.AutoMigrate(&models.Users{}, &models.Cars{}, &models.Customers{}, &models.Garages{}, &models.Orders{})
}

func connectDB(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	// Connect to the database
	DB_URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	var err error
	DB, err = gorm.Open(mysql.Open(DB_URL), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
}

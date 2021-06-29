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

// .env struct
type AppConfig struct {
	Port       string
	DbDriver   string
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
	DbName     string
	JWTSecret  string
}

var DB *gorm.DB

// Define configuration from .env
var Config AppConfig

// Map authorization level for each role
var Mapping = map[string][]string{
	"user":       {"refresh"},
	"superadmin": {"refresh", "basicAdminPermissions"},
}

// It loads values from .env into the system
// Is it the best practice?
// I include mysql configuration on the .env file. Check the .env out!
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}
}

// Connect to the database using DB configuration from .env
// The .env file must be adjusted to comply with your database configuration
func InitDB() {
	connectDB(
		Config.DbDriver,
		Config.DbUser,
		Config.DbPassword,
		Config.DbPort,
		Config.DbHost,
		Config.DbName)
	initMigrate()
}

// Automigrate schema from the struct models
func initMigrate() {
	DB.AutoMigrate(&models.Users{}, &models.Cars{}, &models.Customers{}, &models.Garages{}, &models.Orders{})
}

// Connect to the database
func connectDB(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	DB_URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	var err error
	DB, err = gorm.Open(mysql.Open(DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("FatalError:Cannot connect to DB: \n%v", err)
	}
}

// Get and set configuration from .env
func SetConfig() {
	Config = AppConfig{
		Port:       getEnv("PORT", "8000"),
		DbDriver:   getEnv("DB_DRIVER", ""),
		DbUser:     getEnv("DB_USER", ""),
		DbPassword: getEnv("DB_PASSWORD", ""),
		DbPort:     getEnv("DB_PORT", ""),
		DbHost:     getEnv("DB_HOST", ""),
		DbName:     getEnv("DB_NAME", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}
}

// Prevent a nil environment variable key
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else if defaultVal == "" {
		log.Fatalf("environment variable %s cannot have a nil value", key)
	}
	return defaultVal
}

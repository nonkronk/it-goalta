package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

var validate = validator.New()

type Users struct {
	Id         int       `json:"id" form:"id" gorm:"autoIncrement;not null" validate:"-"`
	Name       string    `json:"name,omitempty" form:"name" gorm:"not null" validate:"required,min=3,max=15"`
	Email      string    `json:"email,omitempty" form:"email" gorm:"not null" validate:"required,email"`
	Password   string    `json:"password,omitempty" form:"password" gorm:"not null" validate:"required,min=3,max=15"`
	Created_at time.Time `json:"created_at" form:"created_at" gorm:"autoCreateTime"`
	Role       string    `json:"role" form:"role" validate:"required,oneof=user admin superadmin"`
}

type UserLoginCreds struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=user admin superadmin"`
}

type UsersResponse struct {
	ID    int    `json:"id,omitempty" form:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Role  string `json:"role"`
}

type JWTClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// Count the number of user objects available in database
func (g *Users) CountAllUsers(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Debug().Model(&Users{}).Count(&count).Error
	return count, err
}

// Print in which tag the validation error occurs
func constructValidationError(err error) ValidationError {
	var validationError ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		fmt.Println(err.Tag())
		var fieldError FieldError
		fieldError.FieldName = err.Tag()
		fieldError.ErrorMessage = "Validation Error in " + err.Tag() + ": " + err.Value().(string)
		validationError.Errors = append(validationError.Errors, fieldError)
	}
	return validationError
}

// Validate JSON input when user register
func (user *Users) Validate() (ValidationError, bool) {
	err := validate.Struct(user)
	if err != nil {
		return constructValidationError(err), false
	}
	return ValidationError{}, true
}

// Validate JSON input when user login
func (user *UserLoginCreds) Validate() (ValidationError, bool) {
	err := validate.Struct(user)
	if err != nil {
		return constructValidationError(err), false
	}
	return ValidationError{}, true
}

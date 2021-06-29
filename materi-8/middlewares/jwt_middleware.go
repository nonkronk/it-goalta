package middlewares

import (
	"project/config"
	"project/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func CreateToken(email, role string, id int, isRefreshToken bool) (string, int64, error) {
	var expiry int64
	if isRefreshToken {
		expiry = int64(time.Now().Add(7 * 24 * time.Hour).Unix())
	} else {
		expiry = int64(time.Now().Add(24 * time.Hour).Unix())
	}
	claims := &models.JWTClaims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(id),
			ExpiresAt: expiry,
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(config.Config.JWTSecret))
	if err != nil {
		return "", 0, err
	}
	return token, claims.StandardClaims.ExpiresAt, nil
}

func GetFieldFromToken(field string, c echo.Context) (string, bool) {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	return claims[field].(string), true
}

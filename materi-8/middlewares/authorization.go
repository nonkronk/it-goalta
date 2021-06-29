package middlewares

import (
	"net/http"
	"project/config"
	"project/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// The authorization level config
var codesMapping = config.Mapping

// Check the role authorization level
func apiCodeValidForRole(apicode, role string) bool {
	for key := range codesMapping {
		if key == role {
			for _, allowedCode := range codesMapping[key] {
				if allowedCode == apicode {
					return true
				}
			}
		}
	}
	return false
}

// Get information from the token and check for valid authorization
func AuthorizationMiddleware(apicode string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			token := user.(*jwt.Token)

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, models.Response{Message: "Access Denied"})
			}
			role := claims["role"].(string)

			if apiCodeValidForRole(apicode, role) {
				return next(c)
			}
			return c.JSON(http.StatusUnauthorized, models.Response{Message: "Access Denied"})
		}
	}
}

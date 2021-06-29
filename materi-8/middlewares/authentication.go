package middlewares

import (
	"project/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Leverage default authentication middleware from echo
func AuthenticationMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(config.Config.JWTSecret),
	})
}

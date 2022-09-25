package auth

import (
	"lah-2022/backend/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthedContext struct {
	echo.Context
	ID           string
	Email        string
	RefreshToken string
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return c.JSON(http.StatusUnauthorized, utils.APIError{
				ErrorMessage: "Unauthorized - no auth token",
			})
		}

		claims := ParseJWT(auth)
		if claims == nil {
			return c.JSON(http.StatusUnauthorized, utils.APIError{
				ErrorMessage: "Unauthorized - invalid token",
			})
		}

		cc := &AuthedContext{
			c,
			claims["sub"].(string),
			claims["email"].(string),
			claims["refresh_token"].(string),
		}

		return next(cc)
	}
}

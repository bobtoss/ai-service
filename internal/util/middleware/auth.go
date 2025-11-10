package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const UserIDContextKey = "user_id"

// AuthMiddleware returns middleware that checks Authorization: Bearer <token>
// The parseFn should return userID string and error
func AuthMiddleware(parseFn func(token string) (string, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization header"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization format"})
			}

			token := parts[1]
			userID, err := parseFn(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired token"})
			}

			c.Set(UserIDContextKey, userID)
			return next(c)
		}
	}
}

// FromContext helper
func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(UserIDContextKey)
	if v == nil {
		return "", false
	}
	uid, ok := v.(string)
	return uid, ok
}

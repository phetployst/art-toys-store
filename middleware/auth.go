package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "Role"
)

type ConfigProvider interface {
	GetJwtSecret() string
}

type ConfigWrapper struct {
	*config.Config
}

func (cw *ConfigWrapper) GetJwtSecret() string {
	return cw.Jwt.AccessTokenSecret
}

type middlewareHandler struct {
	config ConfigProvider
}

func NewMiddlewareHandler(config ConfigProvider) *middlewareHandler {
	return &middlewareHandler{config}
}

func (m *middlewareHandler) JwtMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.ErrUnauthorized
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.ErrUnauthorized
		}
		token = parts[1]

		parsedToken, err := jwt.ParseWithClaims(token, &entities.JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.GetJwtSecret()), nil
		})
		if err != nil || !parsedToken.Valid {
			return echo.ErrUnauthorized
		}

		claims, ok := parsedToken.Claims.(*entities.JwtCustomClaims)
		if !ok {
			return echo.ErrUnauthorized
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)

		return next(c)
	}
}

func (m *middlewareHandler) RbacMiddleware(next echo.HandlerFunc, expectedRole string) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get(ContextRoleKey).(string)
		if !ok || role == "" {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied: Role not found")
		}

		if role != expectedRole {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied: Insufficient permissions")
		}

		return next(c)
	}
}

func (m *middlewareHandler) UserIdParamValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		tokenUserID, ok := c.Get(ContextUserIDKey).(uint)
		if !ok || tokenUserID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID in token")
		}

		paramUserID := c.Param("user_id")
		if paramUserID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "User ID parameter is required")
		}

		paramUserIDUint, err := strconv.ParseUint(paramUserID, 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid User ID format")
		}

		if tokenUserID != uint(paramUserIDUint) {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied: User ID mismatch")
		}

		return next(c)
	}
}

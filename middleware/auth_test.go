package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
)

type MockConfigProvider struct {
	JwtSecret string
}

func (m *MockConfigProvider) GetJwtSecret() string {
	return m.JwtSecret
}
func TestJwtMiddleWare(t *testing.T) {
	t.Run("should pass when token is valid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockConfig := &MockConfigProvider{JwtSecret: "test-secret"}
		handler := NewMiddlewareHandler(mockConfig)

		claims := &entities.JwtCustomClaims{
			UserID: uint(12),
			Role:   "user",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte(mockConfig.GetJwtSecret()))
		req.Header.Set("Authorization", "Bearer "+signedToken)

		middlewareFunc := handler.JwtMiddleWare(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
		err := middlewareFunc(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, uint(12), c.Get(ContextUserIDKey))
		assert.Equal(t, "user", c.Get(ContextRoleKey))
	})

	t.Run("should return unauthorized when token is invalid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockConfig := &MockConfigProvider{JwtSecret: "test-secret"}
		handler := NewMiddlewareHandler(mockConfig)
		req.Header.Set("Authorization", "Bearer invalid-token")

		middlewareFunc := handler.JwtMiddleWare(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
		err := middlewareFunc(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
	})
}

func TestRbacMiddleware(t *testing.T) {
	t.Run("should pass when user has valid role", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockConfig := &MockConfigProvider{}
		handler := NewMiddlewareHandler(mockConfig)

		c.Set(ContextRoleKey, "admin")

		middlewareFunc := handler.RbacMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}, "admin")
		err := middlewareFunc(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return forbidden when user has invalid role", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockConfig := &MockConfigProvider{}
		handler := NewMiddlewareHandler(mockConfig)

		c.Set(ContextRoleKey, "user")

		middlewareFunc := handler.RbacMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}, "admin")
		err := middlewareFunc(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(*echo.HTTPError).Code)
	})
}

func TestUserIdParamValidation(t *testing.T) {
	t.Run("should pass when user ID matches route parameter", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("userID")
		c.SetParamValues("18")

		mockConfig := &MockConfigProvider{}
		handler := NewMiddlewareHandler(mockConfig)

		c.Set(ContextUserIDKey, "18")

		middlewareFunc := handler.UserIdParamValidation(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
		err := middlewareFunc(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return forbidden when user ID does not match route parameter", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("userID")
		c.SetParamValues("21")

		mockConfig := &MockConfigProvider{}
		handler := NewMiddlewareHandler(mockConfig)

		c.Set(ContextUserIDKey, "31")

		middlewareFunc := handler.UserIdParamValidation(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
		err := middlewareFunc(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(*echo.HTTPError).Code)
	})
}

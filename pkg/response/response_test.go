package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	e := echo.New()
	defer e.Close()

	response := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), response)

	message := "An error occurred"
	statusCode := http.StatusBadRequest

	err := ErrorResponse(ctx, statusCode, message)

	assert.NoError(t, err)
	assert.Equal(t, statusCode, response.Code)
	assert.JSONEq(t, `{"message":"An error occurred"}`, response.Body.String())
}

func TestSuccessResponse(t *testing.T) {
	e := echo.New()
	defer e.Close()

	response := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), response)

	data := map[string]interface{}{
		"id":       1,
		"username": "phetploy",
	}
	statusCode := http.StatusOK

	err := SuccessResponse(ctx, statusCode, data)

	assert.NoError(t, err)
	assert.Equal(t, statusCode, response.Code)
	assert.JSONEq(t, `{"id":1,"username":"phetploy"}`, response.Body.String())
}

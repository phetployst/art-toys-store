package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Age      int    `json:"age" validate:"gte=0"`
}

func TestContextWrapper_Bind(t *testing.T) {

	setupContext := func(data TestData) echo.Context {
		requestBody, _ := json.Marshal(data)

		e := echo.New()
		defer e.Close()

		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		return e.NewContext(request, response)
	}

	t.Run("bind input given valid data", func(t *testing.T) {
		data := TestData{Email: "phetploy@example.com", Username: "phetploy", Age: 27}
		ctx := setupContext(data)

		c := &contextWrapper{
			EchoContext: ctx,
			validator:   validator.New(),
		}

		var result TestData
		err := c.Bind(&result)

		assert.NoError(t, err)
		assert.Equal(t, data.Email, result.Email)
		assert.Equal(t, data.Username, result.Username)
		assert.Equal(t, data.Age, result.Age)
	})

	t.Run("bind input given invalid email", func(t *testing.T) {
		data := TestData{Email: "invalid-email", Username: "phetploy", Age: 27}
		ctx := setupContext(data)

		c := &contextWrapper{
			EchoContext: ctx,
			validator:   validator.New(),
		}

		var result TestData
		err := c.Bind(&result)

		assert.Error(t, err)
	})

	t.Run("bind input given invalid username", func(t *testing.T) {
		data := TestData{Email: "phetploy@example.com", Username: "", Age: 27}
		ctx := setupContext(data)

		c := &contextWrapper{
			EchoContext: ctx,
			validator:   validator.New(),
		}

		var result TestData
		err := c.Bind(&result)

		assert.Error(t, err)
	})

	t.Run("bind input given invalid age", func(t *testing.T) {
		data := TestData{Email: "phetploy@example.com", Username: "phetploy", Age: -3}
		ctx := setupContext(data)

		c := &contextWrapper{
			EchoContext: ctx,
			validator:   validator.New(),
		}

		var result TestData
		err := c.Bind(&result)

		assert.Error(t, err)
	})

	t.Run("bind input given invalid data", func(t *testing.T) {
		data := TestData{Email: "", Username: "", Age: -9}
		ctx := setupContext(data)

		c := &contextWrapper{
			EchoContext: ctx,
			validator:   validator.New(),
		}

		var result TestData
		err := c.Bind(&result)

		assert.Error(t, err)
	})

}

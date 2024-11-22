package adapters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	t.Run("get user profile successfully", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("GetUserProfile", "14").Return(&entities.UserProfileResponse{
			UserID:    14,
			Username:  "phetploy",
			FirstName: "Phet",
			LastName:  "Ploy",
			Email:     "phetploy@example.com",
			Address: entities.Address{
				Street:     "123 Green Lane",
				City:       "Bangkok",
				State:      "Central",
				PostalCode: "10110",
				Country:    "Thailand",
			},
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
		}, nil)

		request := httptest.NewRequest(http.MethodGet, "/profile/14", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetParamNames("user_id")
		c.SetParamValues("14")

		err := handler.GetProfile(c)

		expectedResponse := `{
			"user_id": 14,
			"username": "phetploy",
			"first_name": "Phet",
			"last_name": "Ploy",
			"email": "phetploy@example.com",
			"address": {
				"street": "123 Green Lane",
				"city": "Bangkok",
				"state": "Central",
				"postal_code": "10110",
				"country": "Thailand"
			},
			"profile_picture_url": "https://example.com/profiles/14.jpg"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("user credential not found", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("GetUserProfile", "13").Return((*entities.UserProfileResponse)(nil), errors.New("credential not found"))

		request := httptest.NewRequest(http.MethodGet, "/profile/13", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetParamNames("user_id")
		c.SetParamValues("13")

		err := handler.GetProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, `{"message":"User credential not found"}`, response.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("GetUserProfile", "12").Return((*entities.UserProfileResponse)(nil), errors.New("some unexpected error"))

		request := httptest.NewRequest(http.MethodGet, "/profile/12", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetParamNames("user_id")
		c.SetParamValues("12")

		err := handler.GetProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, `{"message":"Internal server error"}`, response.Body.String())
	})
}

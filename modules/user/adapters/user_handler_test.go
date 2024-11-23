package adapters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProfile_user(t *testing.T) {
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

		err := handler.GetUserProfile(c)

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

		err := handler.GetUserProfile(c)

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

		err := handler.GetUserProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, `{"message":"Internal server error"}`, response.Body.String())
	})
}

func TestUpdateUserProfile_user(t *testing.T) {
	t.Run("get update user profile successfully", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("UpdateUserProfile", mock.AnythingOfType("*entities.UserProfile")).Return(&entities.UserProfileResponse{
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

		body := `{"user_id":14,"username":"phetploy","first_name":"Phet","last_name":"Ploy","email":"phetploy@example.com",
		"address":{"street":"123 Green Lane","city":"Bangkok","state":"Central","postal_code":"10110","country":"Thailand"},"profile_picture_url":"https://example.com/profiles/14.jpg"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.UpdateUserProfile(c)

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
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("returns bad request on bind error", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`1234`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.UpdateUserProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"message": "Invalid request data"}`, response.Body.String())
	})

	t.Run("returns bad request on validation error", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		body := `{"user_id":14,"username":"","first_name":"Phet","last_name":"Ploy","email":"invalid-email"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.UpdateUserProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"message": "request data validation failed"}`, response.Body.String())
	})

	t.Run("returns internal server error on use case failure", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("UpdateUserProfile", mock.AnythingOfType("*entities.UserProfile")).Return((*entities.UserProfileResponse)(nil), errors.New("use case error"))

		body := `{"user_id":14,"username":"phetploy","first_name":"Phet","last_name":"Ploy","email":"phetploy@example.com",
		"address":{"street":"123 Green Lane","city":"Bangkok","state":"Central","postal_code":"10110","country":"Thailand"},"profile_picture_url":"https://example.com/profiles/14.jpg"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.UpdateUserProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetAllUserProfile_user(t *testing.T) {
	t.Run("get all user profiles successfully", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockProfiles := []entities.UserProfileResponse{
			{
				UserID: 31, Username: "phetploy", FirstName: "Phet", LastName: "Ploy",
				Email: "phetploy@example.com", Address: entities.Address{
					Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand",
				}, ProfilePictureURL: "https://example.com/profiles/31.jpg",
			},
			{
				UserID: 32, Username: "tonytonychopper", FirstName: "Tony", LastName: "Chopper",
				Email: "tonychopper@example.com", Address: entities.Address{
					Street: "456 Blue Street", City: "Chiang Mai", State: "North", PostalCode: "50200", Country: "Thailand",
				}, ProfilePictureURL: "https://example.com/profiles/32.jpg",
			},
		}

		mockUsecase.On("GetAllUserProfile").Return(int64(2), mockProfiles, nil)

		request := httptest.NewRequest(http.MethodGet, "/user-profiles", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllUserProfile(c)
		expectedResponse := `{
			"count": 2,
			"userProfiles": [
				{"user_id": 31, "username": "phetploy", "first_name": "Phet", "last_name": "Ploy", "email": "phetploy@example.com",
				 "address": {"street": "123 Green Lane", "city": "Bangkok", "state": "Central", "postal_code": "10110", "country": "Thailand"},
				 "profile_picture_url": "https://example.com/profiles/31.jpg"},
				{"user_id": 32, "username": "tonytonychopper", "first_name": "Tony", "last_name": "Chopper", "email": "tonychopper@example.com",
				 "address": {"street": "456 Blue Street", "city": "Chiang Mai", "state": "North", "postal_code": "50200", "country": "Thailand"},
				 "profile_picture_url": "https://example.com/profiles/32.jpg"}
			]
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("GetAllUserProfile").Return(int64(0), ([]entities.UserProfileResponse)(nil), errors.New("unexpected error"))

		request := httptest.NewRequest(http.MethodGet, "/user-profiles", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllUserProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, `{"message":"internal server error"}`, response.Body.String())
	})

	t.Run("no user profiles found", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockUsecase}

		e := echo.New()
		defer e.Close()

		mockUsecase.On("GetAllUserProfile").Return(int64(0), ([]entities.UserProfileResponse)(nil), errors.New("no user profiles found"))

		request := httptest.NewRequest(http.MethodGet, "/user-profiles", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllUserProfile(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, `{"message":"no user profiles found"}`, response.Body.String())
	})
}

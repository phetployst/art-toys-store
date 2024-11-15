package adapters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler_auth(t *testing.T) {

	t.Run("register given valid user should be successful", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewUser", mock.AnythingOfType("*entities.User")).Return(&entities.UserAccount{
			UserID:   uint(1),
			Username: "phetploy",
			Email:    "phetploy@example.com",
		}, nil)

		body := `{"username": "phetploy","email": "phetploy@example.com","password": "12345678","role": "user"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Register(c)

		expectedResponse := `{"user_id":1,"username":"phetploy","email":"phetploy@example.com"}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("register given error during user binding", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{1234}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("register given invalid user", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username": "phetploy"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("register given error email or username already exists", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewUser", mock.AnythingOfType("*entities.User")).Return((*entities.UserAccount)(nil), errors.New("email or username already exists"))

		body := `{"username": "phetploy","email": "phetploy@example.com","password": "12345678","role": "user"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)

	})

	t.Run("register given create user usecase internal server error", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewUser", mock.AnythingOfType("*entities.User")).Return((*entities.UserAccount)(nil), errors.New("internal server error"))

		body := `{"username": "phetploy","email": "phetploy@example.com","password": "12345678","role": "user"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}

func TestLoginHandler_auth(t *testing.T) {
	t.Run("login successful", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("Login", mock.AnythingOfType("*entities.Login"), mock.AnythingOfType("*config.Config")).
			Return(&entities.UserCredential{
				UserID:      1,
				Username:    "phetploy",
				Email:       "phetploy@example.com",
				AccessToken: "access_token",
			}, nil)

		body := `{"username": "phetploy", "password": "password1234"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		expectedResponse := `{"access_token":"access_token", "email":"phetploy@example.com", "user_id":1, "username":"phetploy"}`

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("login with invalid request data", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		body := `{1234}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"message":"Invalid request data"}`, response.Body.String())
	})

	t.Run("login user not found", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("Login", mock.AnythingOfType("*entities.Login"), mock.AnythingOfType("*config.Config")).
			Return((*entities.UserCredential)(nil), errors.New("user not found"))

		body := `{"username": "phetploy", "password": "password1234"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, `{"message":"User not found"}`, response.Body.String())
	})

	t.Run("login invalid password", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("Login", mock.AnythingOfType("*entities.Login"), mock.AnythingOfType("*config.Config")).
			Return((*entities.UserCredential)(nil), errors.New("invalid password"))

		body := `{"username": "phetploy", "password": "wrongpassword"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.JSONEq(t, `{"message":"Invalid password"}`, response.Body.String())
	})

	t.Run("login internal server error", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("Login", mock.AnythingOfType("*entities.Login"), mock.AnythingOfType("*config.Config")).
			Return((*entities.UserCredential)(nil), errors.New("some internal error"))

		body := `{"username": "phetploy", "password": "password1234"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, `{"message":"Internal server error"}`, response.Body.String())
	})

}

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateNewUser(user *entities.User) (*entities.UserAccount, error) {
	args := m.Called(user)
	return args.Get(0).(*entities.UserAccount), args.Error(1)
}

func (m *MockUserUsecase) Login(loginRequest *entities.Login, config *config.Config) (*entities.UserCredential, error) {
	args := m.Called(loginRequest, config)
	return args.Get(0).(*entities.UserCredential), args.Error(1)
}

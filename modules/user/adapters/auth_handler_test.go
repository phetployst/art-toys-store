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

func TestRegisterHandler_auth(t *testing.T) {

	t.Run("register given valid user should be successful", func(t *testing.T) {
		mockService := new(MockUserUsecase)
		handler := &httpUserHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewUser", mock.AnythingOfType("*entities.User")).Return(nil)

		body := `{"username": "phetploy","email": "phetploy@example.com","password": "12345678","role": "user"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
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

		mockService.On("CreateNewUser", mock.AnythingOfType("*entities.User")).Return(errors.New("email or username already exists"))

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

		mockService.On("CreateNewUser", mock.AnythingOfType("*entities.User")).Return(errors.New("internal server error"))

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

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateNewUser(user *entities.User) (*entities.UserAccount, error) {
	args := m.Called(user)
	return nil, args.Error(0)
}

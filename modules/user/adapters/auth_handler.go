package adapters

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/phetployst/art-toys-store/pkg/request"
	"github.com/phetployst/art-toys-store/pkg/response"
)

func (h *httpUserHandler) Register(c echo.Context) error {
	var user entities.User

	wrapper := request.ContextWrapper(c)
	if err := wrapper.Bind(&user); err != nil { // bind and validate
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userAccount, err := h.usecase.CreateNewUser(&user)
	if err != nil {
		if err.Error() == "email or username already exists" {
			return response.ErrorResponse(c, http.StatusConflict, err.Error())
		}
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, userAccount)
}

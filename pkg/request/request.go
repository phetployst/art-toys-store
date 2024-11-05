package request

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		EchoContext echo.Context
		validator   *validator.Validate
	}
)

func ContextWrapper(ctx echo.Context) contextWrapperService {
	return &contextWrapper{
		EchoContext: ctx,
		validator:   validator.New(),
	}
}

func (c *contextWrapper) Bind(data any) error {
	if err := c.EchoContext.Bind(data); err != nil {
		log.Printf("binding error: %v", err)
		return errors.New("invalid input data")
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("validation error: %v", err)
		return errors.New("request data validation failed")
	}

	return nil
}

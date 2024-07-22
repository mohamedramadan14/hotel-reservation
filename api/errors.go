package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var apiError Error
	if errors.As(err, &apiError) {
		return ctx.Status(apiError.Code).JSON(apiError)
	}
	apiError = NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(apiError.Code).JSON(apiError.Err)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Err:  msg,
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid Id given",
	}
}

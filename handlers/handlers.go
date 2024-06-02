package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/types"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h UserHandler) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	user, err := h.userStore.GetUserByID(context.TODO(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}

func (h UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userStore.GetAllUsers(context.TODO())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, users)
}

func (h UserHandler) PostUser(ctx echo.Context) error {
	params, err := DecodeBody[types.CreateUserParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); errs != nil {
		errStrs := ErrorsToStrings(errs)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errStrs})
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	err = h.userStore.CreateUser(context.TODO(), user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}

func DecodeBody[T any](r io.Reader) (T, error) {
	var content T
	err := json.NewDecoder(r).Decode(&content)
	return content, err
}

func ErrorsToStrings(errors []error) []string {
	strs := []string{}
	for _, err := range errors {
		strs = append(strs, err.Error())
	}
	return strs
}

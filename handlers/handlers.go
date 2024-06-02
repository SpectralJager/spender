package handlers

import (
	"context"
	"net/http"

	"github.com/SpectralJager/spender/db"
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
		return err
	}
	return ctx.JSON(http.StatusOK, user)
}

func (h UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userStore.GetAllUsers(context.TODO())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, users)
}

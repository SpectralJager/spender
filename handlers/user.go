package handlers

import (
	"context"
	"net/http"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/middleware"
	"github.com/SpectralJager/spender/types"
	"github.com/SpectralJager/spender/utils"
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
	id := middleware.GetUserIDFromRequest(ctx.Request())
	user, err := h.userStore.GetByID(context.TODO(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"user": user})
}

func (h UserHandler) DeleteUser(ctx echo.Context) error {
	id := middleware.GetUserIDFromRequest(ctx.Request())
	if err := h.userStore.Delete(context.TODO(), id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h UserHandler) PutUser(ctx echo.Context) error {
	params, err := utils.DecodeBody[types.UpdateUserParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	id := middleware.GetUserIDFromRequest(ctx.Request())
	err = h.userStore.Update(context.TODO(), id, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

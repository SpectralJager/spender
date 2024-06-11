package handlers

import (
	"context"
	"net/http"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/middleware"
	"github.com/SpectralJager/spender/types"
	"github.com/SpectralJager/spender/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h AuthHandler) Authenticate(ctx echo.Context) error {
	credentials, err := utils.DecodeBody[types.AuthCredentials](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	user, err := h.userStore.GetByEmail(context.TODO(), credentials.Email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(credentials.Password))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	tokenStr, err := middleware.NewJWTTokenString(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	ctx.Response().Header().Set("X-Api-Token", tokenStr)

	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done"})
}

func (h UserHandler) Register(ctx echo.Context) error {
	params, err := utils.DecodeBody[types.CreateUserParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	_, err = h.userStore.GetByEmail(context.TODO(), params.Email)
	if err == nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "such email already in use"})
	}

	userID, err := h.userStore.Create(context.TODO(), user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	tokenStr, err := middleware.NewJWTTokenString(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	ctx.Response().Header().Set("X-Api-Token", tokenStr)

	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done"})
}

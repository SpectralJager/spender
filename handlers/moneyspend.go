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

type MoneyspendHandler struct {
	moneyspendStore db.SpendStor[types.Moneyspend]
}

func NewMoneyspendHandler(moneyspendStore db.SpendStor[types.Moneyspend]) *MoneyspendHandler {
	return &MoneyspendHandler{
		moneyspendStore: moneyspendStore,
	}
}

func (h MoneyspendHandler) GetAllMonies(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	monies, err := h.moneyspendStore.GetAll(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"monies": monies})
}

func (h MoneyspendHandler) PostMoneyspend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	params, err := utils.DecodeBody[types.CreateMoneyspendParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	moneyspend := types.NewMoneyspendFromParams(params)
	moneyspend.OwnerID = ownerID
	id, err := h.moneyspendStore.Create(context.TODO(), moneyspend)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h MoneyspendHandler) GetMoneyspend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	moneyspend, err := h.moneyspendStore.GetByID(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"money": moneyspend})
}

func (h MoneyspendHandler) PutMoneyspend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	id := ctx.Param("id")
	params, err := utils.DecodeBody[types.UpdateMoneyspendParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	err = h.moneyspendStore.Update(c, id, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h MoneyspendHandler) DeleteMoneyspend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	err := h.moneyspendStore.Delete(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

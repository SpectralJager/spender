package handlers

import (
	"context"
	"net/http"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/types"
	"github.com/SpectralJager/spender/utils"
	"github.com/labstack/echo/v4"
)

type ctxKey string

const (
	ownerIDKey = ctxKey("ownerID")
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
	user, err := h.userStore.GetByID(context.TODO(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"user": user})
}

func (h UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userStore.GetAll(context.TODO())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"users": users})
}

func (h UserHandler) PostUser(ctx echo.Context) error {
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
	userID, err := h.userStore.Create(context.TODO(), user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": userID})
}

func (h UserHandler) DeleteUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if err := h.userStore.Delete(context.TODO(), userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": userID})
}

func (h UserHandler) PutUser(ctx echo.Context) error {
	params, err := utils.DecodeBody[types.UpdateUserParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	userID := ctx.Param("id")
	err = h.userStore.Update(context.TODO(), userID, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": userID})
}

type TimespendHandler struct {
	timespendStore db.SpendStor[types.Timespend]
}

func NewTimespendHandler(timespendStore db.SpendStor[types.Timespend]) *TimespendHandler {
	return &TimespendHandler{
		timespendStore: timespendStore,
	}
}

func (h TimespendHandler) GetAllTimes(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	times, err := h.timespendStore.GetAll(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"timespends": times})
}

func (h TimespendHandler) PostTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	params, err := utils.DecodeBody[types.CreateTimespendParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	timespend := types.NewTimespendFromParams(params)
	timespend.OwnerID = ownerID
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	id, err := h.timespendStore.Create(c, timespend)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h TimespendHandler) GetTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	timespend, err := h.timespendStore.GetByID(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"timespend": timespend})
}

func (h TimespendHandler) PutTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	params, err := utils.DecodeBody[types.UpdateTimespendParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	err = h.timespendStore.Update(c, id, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h TimespendHandler) DeleteTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	err := h.timespendStore.Delete(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})

}

type MoneyspendHandler struct {
	moneyspendStore db.SpendStor[types.Moneyspend]
}

func NewMoneyspendHandler(moneyspendStore db.SpendStor[types.Moneyspend]) *MoneyspendHandler {
	return &MoneyspendHandler{
		moneyspendStore: moneyspendStore,
	}
}

func (h MoneyspendHandler) GetAllMonies(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	monies, err := h.moneyspendStore.GetAll(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"monies": monies})
}

func (h MoneyspendHandler) PostMoneyspend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
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
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	moneyspend, err := h.moneyspendStore.GetByID(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"money": moneyspend})
}

func (h MoneyspendHandler) PutMoneyspend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
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
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	err := h.moneyspendStore.Delete(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})

}

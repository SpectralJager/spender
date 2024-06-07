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

func DecodeBody[T any](r io.Reader) (T, error) {
	var content T
	err := json.NewDecoder(r).Decode(&content)
	return content, err
}

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
	return ctx.JSON(http.StatusOK, echo.Map{"user": user})
}

func (h UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userStore.GetAllUsers(context.TODO())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"users": users})
}

func (h UserHandler) PostUser(ctx echo.Context) error {
	params, err := DecodeBody[types.CreateUserParams](ctx.Request().Body)
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
	userID, err := h.userStore.CreateUser(context.TODO(), user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": userID})
}

func (h UserHandler) DeleteUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if err := h.userStore.DeleteUser(context.TODO(), userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": userID})
}

func (h UserHandler) PutUser(ctx echo.Context) error {
	params, err := DecodeBody[types.UpdateUserParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	userID := ctx.Param("id")
	err = h.userStore.UpdateUser(context.TODO(), userID, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": userID})
}

type TimespendHandler struct {
	ts db.TimespendStore
}

func NewTimespendHandler(store db.TimespendStore) *TimespendHandler {
	return &TimespendHandler{
		ts: store,
	}
}

func (h TimespendHandler) GetAllTimes(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	times, err := h.ts.GetAllTimes(context.TODO(), ownerID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"timespends": times})
}

func (h TimespendHandler) PostTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	params, err := DecodeBody[types.CreateTimespendParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	timespend := types.NewTimespendFromParams(params)
	timespend.OwnerID = ownerID
	id, err := h.ts.CreateTimespend(context.TODO(), timespend)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h TimespendHandler) GetTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	timespend, err := h.ts.GetTimespendByID(context.TODO(), ownerID, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"timespend": timespend})
}

func (h TimespendHandler) PutTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	params, err := DecodeBody[types.UpdateTimespendParams](ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if errs := params.Validate(); len(errs) != 0 {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"validationErrors": errs})
	}
	err = h.ts.UpdateTimespend(context.TODO(), ownerID, id, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})
}

func (h TimespendHandler) DeleteTimespend(ctx echo.Context) error {
	ownerID := ctx.Request().Header.Get("ownerid")
	id := ctx.Param("id")
	err := h.ts.DeleteTimespend(context.TODO(), ownerID, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})

}

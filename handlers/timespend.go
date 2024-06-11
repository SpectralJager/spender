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

type TimespendHandler struct {
	timespendStore db.SpendStor[types.Timespend]
}

func NewTimespendHandler(timespendStore db.SpendStor[types.Timespend]) *TimespendHandler {
	return &TimespendHandler{
		timespendStore: timespendStore,
	}
}

func (h TimespendHandler) GetAllTimes(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	times, err := h.timespendStore.GetAll(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"timespends": times})
}

func (h TimespendHandler) PostTimespend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
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
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	timespend, err := h.timespendStore.GetByID(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"timespend": timespend})
}

func (h TimespendHandler) PutTimespend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
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
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	id := ctx.Param("id")
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)
	err := h.timespendStore.Delete(c, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"resutl": "done", "id": id})

}

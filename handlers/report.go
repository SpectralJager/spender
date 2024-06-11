package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/middleware"
	"github.com/SpectralJager/spender/types"
	"github.com/SpectralJager/spender/utils"
	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	timespendStore  db.MongoTimespendStore
	moneyspendStore db.MongoMoneyspendStore
}

func NewReportHandler(timespendStore db.MongoTimespendStore, moneyspendStore db.MongoMoneyspendStore) *ReportHandler {
	return &ReportHandler{
		timespendStore:  timespendStore,
		moneyspendStore: moneyspendStore,
	}
}

func (h ReportHandler) GetTotalSpend(ctx echo.Context) error {
	ownerID := middleware.GetUserIDFromRequest(ctx.Request())
	c := context.WithValue(context.Background(), ownerIDKey, ownerID)

	times, err := h.timespendStore.GetAll(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	monies, err := h.moneyspendStore.GetAll(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	start := ctx.QueryParam("start")
	end := ctx.QueryParam("end")
	if len(start) != 0 && len(end) != 0 {
		dateStart, err := time.Parse(time.DateOnly, start)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		dateEnd, err := time.Parse(time.DateOnly, end)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		times = utils.Map(times, func(timespend types.Timespend) bool {
			return timespend.Date.After(dateStart) && timespend.Date.Before(dateEnd)
		})
		monies = utils.Map(monies, func(moneyspend types.Moneyspend) bool {
			return moneyspend.Date.After(dateStart) && moneyspend.Date.Before(dateEnd)
		})
	}

	totalTime := types.Timespend{
		OwnerID: ownerID,
	}
	for _, time := range times {
		totalTime.Duration += time.Duration
	}

	totalMoney := types.Moneyspend{
		OwnerID: ownerID,
	}
	for _, money := range monies {
		totalMoney.Money += money.Money
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"timespend":  totalTime,
		"moneyspend": totalMoney,
	})
}

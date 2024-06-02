package handlers

import (
	"net/http"

	"github.com/SpectralJager/spender/types"
	"github.com/labstack/echo/v4"
)

func HandleGetUsers(ctx echo.Context) error {
	users := []types.User{
		{FirstName: "James", LastName: "Smith"},
		{FirstName: "Charls", LastName: "Jonhson"},
	}
	return ctx.JSON(http.StatusOK, users)
}

func HandleGetUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"user": "James",
	})
}

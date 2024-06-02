package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/SpectralJager/spender/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	listenAddr := flag.String("addr", ":8080", "the listhen addres of api server")
	flag.Parse()

	app := echo.New()

	app.GET("/", handleIndex)

	apiv1 := app.Group("/api/v1")
	apiv1.GET("/user", handlers.HandleGetUsers)
	apiv1.GET("/user/:id", handlers.HandleGetUser)

	if err := app.Start(*listenAddr); err != nil {
		log.Fatalf("something goes wrong -> %v", err)
	}
}

func handleIndex(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, world")
}

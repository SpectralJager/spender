package main

import (
	"context"
	"flag"
	"log"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBURI    = "mongodb://localhost:27017"
	DBNAME   = "spender"
	USERCOLL = "users"
)

func main() {
	listenAddr := flag.String("addr", ":8080", "the listhen addres of api server")
	flag.Parse()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DBURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	userHandler := handlers.NewUserHandler(db.NewMongoUserStore(client, DBNAME, USERCOLL))

	app := echo.New()
	apiv1 := app.Group("/api/v1", middleware.Logger())
	apiv1.GET("/user", userHandler.GetUsers)
	apiv1.POST("/user", userHandler.PostUser)
	apiv1.GET("/user/:id", userHandler.GetUser)
	apiv1.PUT("/user/:id", userHandler.PutUser)
	apiv1.DELETE("/user/:id", userHandler.DeleteUser)

	// userHandler.InitHandlers(apiv1)

	if err := app.Start(*listenAddr); err != nil {
		log.Fatalf("something goes wrong -> %v", err)
	}
}

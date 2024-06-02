package main

import (
	"context"
	"flag"
	"log"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/handlers"
	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "spender"
const userColl = "users"

func main() {
	listenAddr := flag.String("addr", ":8080", "the listhen addres of api server")
	flag.Parse()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	userHandler := handlers.NewUserHandler(db.NewMongoUserStore(client))

	app := echo.New()
	apiv1 := app.Group("/api/v1")
	apiv1.GET("/user", userHandler.GetUsers)
	apiv1.GET("/user/:id", userHandler.GetUser)

	// userHandler.InitHandlers(apiv1)

	if err := app.Start(*listenAddr); err != nil {
		log.Fatalf("something goes wrong -> %v", err)
	}
}

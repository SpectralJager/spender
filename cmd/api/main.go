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
	DBURI          = "mongodb://localhost:27017"
	DBNAME         = "spender"
	USERCOLL       = "users"
	TIMESPENDCOLL  = "timespends"
	MONEYSPENDCOLL = "moneyspends"
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

	userStore := db.NewMongoUserStore(client, DBNAME, USERCOLL)
	timespendStore := db.NewMongoTimespendStore(client, DBNAME, TIMESPENDCOLL)
	moneyspendStore := db.NewMongoMoneyspendStore(client, DBNAME, MONEYSPENDCOLL)

	userHandler := handlers.NewUserHandler(userStore)
	timespendHandler := handlers.NewTimespendHandler(timespendStore)
	moneyspendHandler := handlers.NewMoneyspendHandler(moneyspendStore)
	reportHandler := handlers.NewReportHandler(timespendStore, moneyspendStore)

	app := echo.New()
	apiv1 := app.Group("/api/v1", middleware.Logger())
	// User api
	apiv1.GET("/user", userHandler.GetUsers)
	apiv1.POST("/user", userHandler.PostUser)
	apiv1.GET("/user/:id", userHandler.GetUser)
	apiv1.PUT("/user/:id", userHandler.PutUser)
	apiv1.DELETE("/user/:id", userHandler.DeleteUser)
	// Timespend api
	apiv1.GET("/timespend", timespendHandler.GetAllTimes)
	apiv1.POST("/timespend", timespendHandler.PostTimespend)
	apiv1.GET("/timespend/:id", timespendHandler.GetTimespend)
	apiv1.PUT("/timespend/:id", timespendHandler.PutTimespend)
	apiv1.DELETE("/timespend/:id", timespendHandler.DeleteTimespend)
	// Moneyspend api
	apiv1.GET("/moneyspend", moneyspendHandler.GetAllMonies)
	apiv1.POST("/moneyspend", moneyspendHandler.PostMoneyspend)
	apiv1.GET("/moneyspend/:id", moneyspendHandler.GetMoneyspend)
	apiv1.PUT("/moneyspend/:id", moneyspendHandler.PutMoneyspend)
	apiv1.DELETE("/moneyspend/:id", moneyspendHandler.DeleteMoneyspend)
	// Report api
	apiv1.GET("/report/total", reportHandler.GetTotalSpend)
	apiv1.GET("/report/moneyspend", reportHandler.GetMoneyspends)
	apiv1.GET("/report/timespend", reportHandler.GetTimespends)

	if err := app.Start(*listenAddr); err != nil {
		log.Fatalf("something goes wrong -> %v", err)
	}
}

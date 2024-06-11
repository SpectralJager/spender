package main

import (
	"context"
	"flag"
	"log"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/handlers"
	"github.com/SpectralJager/spender/middleware"
	"github.com/labstack/echo/v4"

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

	authHandler := handlers.NewAuthHandler(userStore)
	userHandler := handlers.NewUserHandler(userStore)
	timespendHandler := handlers.NewTimespendHandler(timespendStore)
	moneyspendHandler := handlers.NewMoneyspendHandler(moneyspendStore)
	reportHandler := handlers.NewReportHandler(timespendStore, moneyspendStore)

	app := echo.New()
	apiv1 := app.Group("/api/v1")
	// Authentication api
	authApi := apiv1.Group("/auth")
	authApi.POST("/login", authHandler.Authenticate)
	authApi.POST("/register", userHandler.Register)
	// User api
	userApi := apiv1.Group("/user", middleware.JWTAuthentication)
	userApi.GET("", userHandler.GetUser)
	userApi.PUT("", userHandler.PutUser)
	userApi.DELETE("", userHandler.DeleteUser)
	// Timespend api
	timespendApi := apiv1.Group("/timespend", middleware.JWTAuthentication)
	timespendApi.GET("", timespendHandler.GetAllTimes)
	timespendApi.POST("", timespendHandler.PostTimespend)
	timespendApi.GET("/:id", timespendHandler.GetTimespend)
	timespendApi.PUT("/:id", timespendHandler.PutTimespend)
	timespendApi.DELETE("/:id", timespendHandler.DeleteTimespend)
	// Moneyspend api
	moneyspendApi := apiv1.Group("/moneyspend", middleware.JWTAuthentication)
	moneyspendApi.GET("", moneyspendHandler.GetAllMonies)
	moneyspendApi.POST("", moneyspendHandler.PostMoneyspend)
	moneyspendApi.GET("/:id", moneyspendHandler.GetMoneyspend)
	moneyspendApi.PUT("/:id", moneyspendHandler.PutMoneyspend)
	moneyspendApi.DELETE("/:id", moneyspendHandler.DeleteMoneyspend)
	// Report api
	reportApi := apiv1.Group("/report", middleware.JWTAuthentication)
	reportApi.GET("/total", reportHandler.GetTotalSpend)

	if err := app.Start(*listenAddr); err != nil {
		log.Fatalf("something goes wrong -> %v", err)
	}
}

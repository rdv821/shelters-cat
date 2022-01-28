package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/catService/docs"
	"github.com/catService/internal/config"
	"github.com/catService/internal/handlers"
	"github.com/catService/internal/repository"
	"github.com/catService/internal/service"
	"github.com/catService/internal/validator"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title        Cats API
// @version      1.0
// @description  API server for shelters cats

// @host      localhost:9090
// @BasePath  /v1/
// @schemes   http

func main() {
	const timeout = 20 * time.Second
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	// config
	cfg, err := config.New()
	if err != nil {
		logrus.Fatalf("Can't initialized config: %+v", err)
	}

	var rps repository.SheltersCatRepository
	client := NewRedis(cfg.RedisURL)
	redisRepository := repository.NewLocalCache(ctx, client)

	switch cfg.DBType {
	case "postgres":
		db := NewPostgresDB(cfg.PostgresURL)
		rps = repository.NewPostgresRepository(db)
	case "mongo":
		db := NewMongoDB(cfg.MongoURL)
		rps = repository.NewMongoRepository(db)
	default:
		logrus.Fatalf("Unknown db type %v", cfg.DBType)
	}

	srv := service.NewService(rps, redisRepository)
	catHandler := handlers.NewCat(srv)

	e := echo.New()
	e.Validator = validator.NewValidator()

	v1 := e.Group("/v1")
	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	catRouters := v1.Group("/cat")
	catRouters.POST("/", catHandler.Create)
	catRouters.GET("/:id", catHandler.Get)
	catRouters.DELETE("/:id", catHandler.Delete)
	catRouters.PUT("/:id", catHandler.Update)

	go func() {
		err = e.Start(cfg.ServerPort)
		if errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("Server stopped: %v", err)
		} else {
			logrus.Fatalf("Can't start server %v", err)
		}
	}()

	<-ctx.Done()
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := e.Shutdown(ctxWithTimeout); err != nil {
		logrus.Fatalf("Can't shutdown server gracefully: %v", err)
	}
}

// NewPostgresDB create connection to db
func NewPostgresDB(dbURL string) *pgxpool.Pool {
	const timeout = 10 * time.Minute
	ctxWithTimeout, cancelFunction := context.WithTimeout(context.Background(), timeout)
	defer cancelFunction()
	pool, err := pgxpool.Connect(ctxWithTimeout, dbURL)
	logrus.Infof("Postgres connection was started...")
	if err != nil {
		logrus.Fatalf("Unable to connection to database: %v", err)
	}

	return pool
}

// NewMongoDB create connection to db
func NewMongoDB(dbURL string) *mongo.Database {
	const timeout = 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	logrus.Infof("Mongo connection was started...")
	if err != nil {
		logrus.Fatalf("Unable to connection to database: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Fatalf("connection was not established: %v", err)
	}

	return client.Database("cats")
}

// NewRedis create connection
func NewRedis(redisURL string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("connection to redisd failed: %v", err)
	}
	logrus.Info("Redis connection was started...")

	return client
}

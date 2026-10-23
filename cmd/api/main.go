package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoURL is read from the MONGO_URL environment variable.
var mongoURL = os.Getenv("MONGO_URL")

const (
	webPort  = "80"
	rpcPort  = "5001"
	gRPCPort = "5001"
)

type Config struct {
	Repo data.Repository
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{}
	app.setupRepo(mongoClient)

	log.Println("starting service on", webPort)
	if mongoURL == "" {
		log.Fatal("MONGO_URL must be set")
	}
	app.Serve()
}

func (app *Config) Serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connceting:", err)
		return nil, err
	}
	log.Println("connected to mongo db")
	return c, nil
}

func (app *Config) setupRepo(mongoClient *mongo.Client) {
	loggerRepo := data.NewLoggerRepository(mongoClient)
	app.Repo = loggerRepo
}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/petrovi4ev/bitmedia-test/internal/apiserver"
	"github.com/petrovi4ev/bitmedia-test/internal/config"
	"github.com/petrovi4ev/bitmedia-test/internal/migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	migrateUp bool
)

func init() {
	flag.BoolVar(&migrateUp, "migrate", false, "a bool")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.New()
	flag.Parse()
	db := getClient(*cfg)
	if migrateUp {
		migrate.Up(cfg.DbName, db)
	}

	server := apiserver.New(db, *cfg)
	server.Start()
}

func getClient(cfg config.Config) *mongo.Client {
	s := fmt.Sprintf("mongodb://%s:%s", cfg.DbHost, cfg.DbPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(s))
	check(err)
	err = client.Connect(context.TODO())
	check(err)
	err = client.Ping(context.TODO(), nil)
	check(err)

	log.Println("Connected to MongoDB!")

	return client
}

func check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

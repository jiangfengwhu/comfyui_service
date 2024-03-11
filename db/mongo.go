package db

import (
	"comfyui_service/utils"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type DB struct {
	client *mongo.Client
	db     *mongo.Database
	task   *mongo.Collection
	user   *mongo.Collection
}

var db = DB{}

func Init() {
	println("connect db")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(utils.Config.MongoUrl))
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database("llm")
	db.client = client
	db.db = database
	db.task = database.Collection("tasks")
	db.user = database.Collection("user")
}
func CloseDB() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

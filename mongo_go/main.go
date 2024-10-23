package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kaung-minkhant/go_projs/mongo_go/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
  client := getClient()
  defer func() {
    if err := client.Disconnect(context.Background()); err != nil {
      panic(err)
    }
  }()

	r := httprouter.New()
	uc := controllers.NewUserController(client)
	r.GET("/users/:id", uc.GetUser)
	r.POST("/users", uc.CreateUser)
	r.DELETE("/users/:id", uc.DeleteUser)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getClient() *mongo.Client {
  client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()
  err := client.Ping(ctx, readpref.Primary())
  if err != nil {
    panic(err)
  }
	return client
}

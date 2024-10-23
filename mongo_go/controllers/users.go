package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kaung-minkhant/go_projs/mongo_go/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserController interface {
	GetUser(http.ResponseWriter, *http.Request, httprouter.Params)
	CreateUser(http.ResponseWriter, *http.Request, httprouter.Params)
	DeleteUser(http.ResponseWriter, *http.Request, httprouter.Params)
}

type UserC struct {
	client *mongo.Client
}

func NewUserController(c *mongo.Client) UserController {
	return &UserC{c}
}

func (uc *UserC) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{}

  err = uc.client.Database("mongo_golang").Collection("users").FindOne(context.TODO(), bson.D{{"_id", oid}}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Get User Marshal Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (uc *UserC) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var newUser models.User
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, &newUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	newUser.Id = bson.NewObjectID()


  _, err := uc.client.Database("mongo_golang").Collection("users").InsertOne(context.TODO(), newUser)
	if err != nil {
		fmt.Println("Create User Insert Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println("Create User Marshal Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (uc *UserC) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
  oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("id is wrong"))
		return
	}
  _, err = uc.client.Database("mongo_golang").Collection("users").DeleteOne(context.TODO(), bson.D{{"_id", oid}})
	if err != nil {
		fmt.Println("Delete user error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Deleted user %s", id)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const DB_NAME = "fiber_hrm"
const MONGO_URL = "mongodb://root:example@localhost:27017"

type Employee struct {
	ID     string  `json:"id,ommitempty" bson:"_id,ommitempty"`
	Name   string  `json:"name"`
	Salary float32 `json:"salary"`
	Age    int32   `json:"age"`
}

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(MONGO_URL))
	if err != nil {
		return err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

  fmt.Println("Db connection successful!")
	db := client.Database(DB_NAME)
	mg.Client = client
	mg.Db = db
	return nil
}

func GetEmployees(c *fiber.Ctx) error {
	var employees []Employee
	collection := mg.Db.Collection("employees")
	query := bson.D{}
	cursor, err := collection.Find(c.Context(), query)
	if err != nil {
		fmt.Println("Getting Employees failed", err)
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Something went wrong"))
	}
	defer cursor.Close(c.Context())
	for cursor.Next(c.Context()) {
		var employee Employee
		if err := cursor.Decode(&employee); err != nil {
			fmt.Println("Decoding Employees in GetEmployees failed", err)
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Something went wrong"))
		}
		employees = append(employees, employee)
	}
	return c.JSON(employees)
}

func GetEmployee(c *fiber.Ctx) error {
  id := c.Params("id")
  _, err := bson.ObjectIDFromHex(id)
  if err != nil {
    fmt.Println("Parsig id in GetEmployee failed", err)
    return c.Status(fiber.StatusBadRequest).SendString("Bad request")
  }
  collection := mg.Db.Collection("employees")
  var employee Employee
  result := collection.FindOne(c.Context(), bson.D{{Key: "_id", Value: id}})
  if err := result.Err(); err != nil {
    fmt.Println("Gettng Employee failed", err)
    switch {
    case errors.Is(err, mongo.ErrNoDocuments):
      return c.Status(fiber.StatusNotFound).SendString("Not found")
    default:
      return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
    }
  }
  if err := result.Decode(&employee); err != nil {
    fmt.Println("Parsing Employee in GetEmployee failed", err)
    return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
  }
	return c.Status(fiber.StatusOK).JSON(employee)
}
func CreateEmployee(c *fiber.Ctx) error {
	employee := new(Employee)
	if err := c.BodyParser(&employee); err != nil {
		fmt.Println("Decoding Body in Create Employee failed", err)
		switch {
		case errors.Is(err, fiber.ErrUnprocessableEntity):
			return c.Status(fiber.StatusBadRequest).Send([]byte("Bad Request"))
		default:
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Something went wrong"))
		}
	}
	collection := mg.Db.Collection("employees")
	employee.ID = bson.NewObjectID().Hex()

	result, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		fmt.Println("Creating new Employee failed", err)
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Something went wrong"))
	}
	id := result.InsertedID

	filter := bson.D{{Key: "_id", Value: id}}
	var insertedEmployee Employee
	collection.FindOne(c.Context(), filter).Decode(&insertedEmployee)
	return c.Status(fiber.StatusCreated).JSON(insertedEmployee)
}
func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := bson.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Decoding id in Update Employee failed", err)
		return c.Status(fiber.StatusNotFound).Send([]byte("Not found"))
	}
	var employee Employee

	if err := c.BodyParser(&employee); err != nil {
		fmt.Println("Decoding body in Update Employee failed", err)
		return c.Status(fiber.StatusBadRequest).Send([]byte("Bad request"))
	}

	colletion := mg.Db.Collection("employees")

	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}

	var updatedEmployee Employee
	result := colletion.FindOneAndUpdate(c.Context(), query, update)
	if err := result.Err(); err != nil {
		fmt.Println("Updating employee failed", err)
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return c.Status(fiber.StatusNotFound).Send([]byte("Employee not found"))
		default:
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Something went wrong"))
		}
	}
	colletion.FindOne(c.Context(), bson.D{{Key: "_id", Value: id}}).Decode(&updatedEmployee)

	return c.Status(fiber.StatusAccepted).JSON(updatedEmployee)
}

func DeleteEmployee(c *fiber.Ctx) error {
  id := c.Params("id")
  _, err := bson.ObjectIDFromHex(id)
  if err != nil {
    fmt.Println("Parsing id in DeleteEmployee failed", err)
    return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
  }
  collection := mg.Db.Collection("employees")
  result, err := collection.DeleteOne(c.Context(), bson.D{{Key: "_id", Value: id}})
  if err != nil {
    fmt.Println("Deleting employee failed", err)
    return c.Status(fiber.StatusInternalServerError).SendString("Someting went wrong")
  }
  if result.DeletedCount < 1 {
    return c.Status(fiber.StatusNotFound).SendString("Not found")
  }
  return c.Status(fiber.StatusAccepted).SendString(fmt.Sprintf("Employee with ID %s deleted\n", id))
}

func SetupRoutes(app *fiber.App) {
	app.Get("/employee", GetEmployees)
	app.Get("/employee/:id", GetEmployee)
	app.Post("/employee", CreateEmployee)
	app.Put("/employee/:id", UpdateEmployee)
	app.Delete("/employee/:id", DeleteEmployee)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Connect(); err != nil {
		log.Fatal("Cannot connect to db", err)
		return
	}
	defer func() {
		if err := mg.Client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()
	SetupRoutes(app)

  log.Fatal(app.Listen(":8080"))
}

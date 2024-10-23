package database

import (
	"database/sql"
  _ "github.com/lib/pq"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func CreateConnection() *sql.DB {
  err := godotenv.Load(".env")
  if err != nil {
    panic("Error opening .env file")
  }

  db, err := sql.Open("postgres", os.Getenv("POSTGRES_CON"))

  if err != nil {
    fmt.Println("Error creating db object", err)
    panic("Error creating db object")
  }
  err = db.Ping()
  if err != nil {
    fmt.Println("Error connecting to db", err)
    panic("Error connecting to db")
  }
  return db
}

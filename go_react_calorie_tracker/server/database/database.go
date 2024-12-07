package database

import (
	"log"
	"os"
)

type Database interface {
  connect() (interface{}, error)
  disconnect() error
  Ping() error
}

func InitDatabase(db Database) {
  err := db.Ping()
  if err != nil {
    log.Printf("%s\n", err)
    os.Exit(1)
  }
}

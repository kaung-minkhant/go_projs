package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
  DBConn *gorm.DB
)

func InitDatabase() {
  db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
  if err != nil {
    log.Fatal("Cannot opent db", err)
    panic("Error")
  }
  fmt.Println("Connection to DB opened")
  DBConn = db
}


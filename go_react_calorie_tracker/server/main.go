package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kaung-minkhant/go_projs/go_react_calorie_tracker/routes"
)

func main() {
  if err := godotenv.Load(".env"); err != nil {
    panic("Cannot load env")
  }
  port := os.Getenv("PORT")
  engine := routes.Init()
  
  engine.Run(":"+port)
}

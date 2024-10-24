package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-io/slacker"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading env:", err)
		return
	}

	bot := slacker.NewClient(os.Getenv("BOT_TOKEN"), os.Getenv("APP_TOKEN"))

	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "my yob is <year>",
		Description: "Give me your year of birth and i will give you your age",
		Examples: []string{
			"my yob is 2019",
			"my yob is 2024",
		},
		Handler: func(botCtx *slacker.CommandContext) {
			year := botCtx.Request().Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				log.Fatal("Cannot convert to yob integer", err)
				r := fmt.Sprintf("Error in calculating age: %s", err)
				botCtx.Response().Reply(r)
				return
			}
			age := time.Now().Year() - yob

			r := fmt.Sprintf("Age is %d", age)

			botCtx.Response().Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bot.Listen(ctx); err != nil {
		log.Fatal("Error listening to bot", err)
		return
	}

}

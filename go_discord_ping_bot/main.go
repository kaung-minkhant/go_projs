package main

import (
	"log"

	"github.com/kaung-minkhant/go_projs/go_discord_ping_bot/bot"
	"github.com/kaung-minkhant/go_projs/go_discord_ping_bot/config"
)


func main() {
  if err := config.LoadConfig("./config.json"); err != nil {
    log.Fatal("Cannot read configuration file", err)
    return
  }

  bot.Start()

  <-make(chan struct{}) // stalling 
  return
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Edw590/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker/v2"
	"github.com/tidwall/gjson"
  witai "github.com/kaung-minkhant/wit-go"
)

type SlackApiConfig struct {
	BotToken string
	AppToken string
}

type WitApiConfig struct {
	Token string
}

type WoframApiConfig struct {
	AppId string
}

var (
	SlackApi   SlackApiConfig
	WitApi     WitApiConfig
	WolframApi WoframApiConfig
)

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Cannot load env: %s\n", err)
		os.Exit(1)
	}

  SlackApi = SlackApiConfig{
    BotToken: os.Getenv("SLACK_BOT_TOKEN"),
    AppToken: os.Getenv("SLACK_APP_TOKEN"),
  }

  WitApi = WitApiConfig{
    Token: os.Getenv("WIT_TOKEN"),
  }

  WolframApi = WoframApiConfig{
    AppId: os.Getenv("WOLFRAM_APP_ID"),
  }
}

func main() {
	loadEnv()

  bot := slacker.NewClient(SlackApi.BotToken, SlackApi.AppToken)

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  witClient := witai.NewClient(WitApi.Token)

  wolframClient := &wolfram.Client{AppID: WolframApi.AppId}
  
  bot.AddCommand(&slacker.CommandDefinition{
    Command: "Hey, <message>",
    Description: "send any question to wolfram",
    Examples: []string{
      "Hey, who is the president of spain",
    },
    Handler: func (c *slacker.CommandContext) {
      query := c.Request().Param("message")
      fmt.Printf("Sending query for %s\n", query)

      // query the wit ai
      messageRequest := &witai.MessageRequest{
        Query: query,
      }
      msg, err := witClient.Parse(messageRequest)
      if err != nil {
        fmt.Printf("Error when querying wit ai: %s\n", err)
        c.Response().Reply("Wit ai error")
        return
      }

      formatedData, _ := json.MarshalIndent(msg, "", "    ")
      roughData := string(formatedData)

      wolframQuery := gjson.Get(roughData, "entities.wit$wolfram_search_query:wolfram_search_query.0.value").String()

      result, err := wolframClient.GetSpokentAnswerQuery(wolframQuery, wolfram.Metric, 1000)
      if err != nil {
        fmt.Printf("Error when querying wolfram: %s\n", err)
        c.Response().Reply("Wolfram error")
        return
      }
      fmt.Println(result)

      c.Response().Reply(result)
    },
  })

  if err := bot.Listen(ctx); err != nil {
    log.Fatal("Error from bot listen", err)
    os.Exit(1)
  }
}

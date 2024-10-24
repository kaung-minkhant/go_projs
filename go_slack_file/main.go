package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
  if err := godotenv.Load(".env"); err != nil {
    panic("Cannot read env file\n")
  }
  api := slack.New(os.Getenv("BOT_TOKEN"))
  channelArr := []string{os.Getenv("CHANNEL_ID")}
  fileArr := []string{"cat.jpg"}

  for i := 0; i < len(fileArr); i++ {
    info, err := os.Stat(fileArr[0])
    if err != nil {
      fmt.Printf("Cannot read file: %v\n", err)
      return
    }
    params := slack.UploadFileV2Parameters{
      Channel: channelArr[0], 
      File: fileArr[i],
      Filename: "cat.jpg",
      FileSize: int(info.Size()),
    }

    file, err := api.UploadFileV2(params)
    if err != nil {
      fmt.Printf("Error: %s\n", err)
      return
    }

    fmt.Printf("Name: %s\n", file.Title)
  }
}

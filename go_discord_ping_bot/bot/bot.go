package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaung-minkhant/go_projs/go_discord_ping_bot/config"
)

var BotID string
var goBot *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
  {
    Name: "echo",
    Description: "Say something through a bot",
    Options: []*discordgo.ApplicationCommandOption{
      {
        Name: "message",
        Description: "Contents of the message",
        Type: discordgo.ApplicationCommandOptionString,
        Required: true,
      },
    },
  },
  {
    Name: "whowelove",
    Description: "Say something through a bot",
  },
}

func Start() {
	goBot, err := discordgo.New(config.Prefix + config.Token)
	if err != nil {
		log.Fatal("Cannot start bot", err)
		return
	}

	user, err := goBot.User("@me")
	if err != nil {
		log.Fatal("Cannot get user", err)
		return
	}

	BotID = user.ID

	goBot.AddHandler(messageHandler)

  goBot.AddHandler(readyHandler)

  goBot.AddHandler(echoHandler)

  goBot.AddHandler(whoWeLoveHandler)

  _, err = goBot.ApplicationCommandBulkOverwrite(config.AppId, config.ServerId, commands)
  if err != nil {
    log.Fatal("Could not regster commands:", err)
    return
  }

	if err := goBot.Open(); err != nil {
		log.Fatal("Cannot open discord connection", err)
		return
	}

	fmt.Println("Yay! PingBot is running!")
}

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func parseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (optionMap) {
  om := make(optionMap)
  for _, opt := range options {
    om[opt.Name] = opt
  }
  return om
}

func interactionAuthor(i *discordgo.Interaction) *discordgo.User {
  if i.Member != nil {
    return i.Member.User
  }
  return i.User
}

func whoWeLoveHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
  if i.Type != discordgo.InteractionApplicationCommand {
    return
  }

  data := i.ApplicationCommandData()
  if data.Name != "whowelove" {
    return
  }

  handleWhoWeLove(s, i, parseOptions(data.Options))
}

func echoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
  if i.Type != discordgo.InteractionApplicationCommand {
    return
  }

  data := i.ApplicationCommandData()
  if data.Name != "echo" {
    return
  }
  
  handleEcho(s, i, parseOptions(data.Options))
}

func handleWhoWeLove(s *discordgo.Session, i *discordgo.InteractionCreate, opts optionMap) {
  builder := new(strings.Builder)
  builder.WriteString("We love shunn")
  err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Content: builder.String(),
    },
  })
  if err != nil {
    log.Fatal("Could not respond to interaction", err)
  }
}

func handleEcho(s *discordgo.Session, i *discordgo.InteractionCreate, opts optionMap) {
  builder := new(strings.Builder)
  // if v, ok := opts["author"]; ok && v.BoolValue() {
  //   author := interactionAuthor(i.Interaction)
  //   builder.WriteString("**" + author.String() + "** says: ")
  // }
  builder.WriteString(opts["message"].StringValue())

  err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Content: builder.String(),
    },
  })

  if err != nil {
    log.Fatal("Could not respond to interaction", err)
  }
}


func readyHandler(s *discordgo.Session, r *discordgo.Ready) {
  log.Printf("Loggd in as %s", r.User.String())
}

func messageHandler(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "ping" {
    _, err := session.ChannelMessageSend(m.ChannelID, "pong", func(cfg *discordgo.RequestConfig) {
			cfg.MaxRestRetries = 3
		})
    if err != nil {
      log.Fatal("Cannot send message", err)
      return
    }
	}
}

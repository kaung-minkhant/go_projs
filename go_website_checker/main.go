package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "HealthChecker",
		Usage: "Tool to check whether a website is up or down.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "domain",
				Aliases:  []string{"d"},
				Usage:    "Domain name to check.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Usage:    "Port number to check",
				Required: false,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
      port := command.String("port")
      if port == "" {
        port = "80"
      }

      status := Check(command.String("domain"), port)
      fmt.Println("Status:", status)
			return nil
		},
	}

  err := cmd.Run(context.Background(), os.Args)
  if err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}

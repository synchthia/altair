package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/altair/bot"
	"gitlab.com/Startail/altair/logger"
	"gitlab.com/Startail/altair/nebula"
)

func main() {
	// Init Logger
	logger.Init()

	// Init
	logrus.Printf("[ALTAIR] Starting ALTAIR Bot...")

	// gRPC
	go func() {
		nebula.NewClient()
	}()

	// Bot
	wait := make(chan struct{})
	go func() {
		defer close(wait)

		discordToken := os.Getenv("DISCORD_TOKEN")
		if len(discordToken) == 0 {
			panic("DISCORD_TOKEN is not defined!")
		}

		err := bot.InitDiscordBot(discordToken)
		if err != nil {
			panic("Discord Bot Error!")
		}
	}()
	<-wait
}

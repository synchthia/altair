package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/altair/bot"
	"gitlab.com/Startail/altair/logger"
	"gitlab.com/Startail/altair/nebula"
	"gitlab.com/Startail/altair/stream"
	"gitlab.com/Startail/altair/systera"
)

func main() {
	// Init Logger
	logger.Init()

	// Init
	logrus.Printf("[ALTAIR] Starting ALTAIR Bot...")

	// Redis
	go func() {
		redisAddr := os.Getenv("REDIS_ADDRESS")
		if len(redisAddr) == 0 {
			redisAddr = "localhost:6379"
		}
		stream.NewRedisPool(redisAddr)

		// Subscribe
		go func() {
			stream.PunishSubs()
		}()

		go func() {
			stream.ReportSubs()
		}()
	}()

	// gRPC
	go func() {
		systera.NewClient()
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

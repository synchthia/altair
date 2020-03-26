package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/altair/bot"
	"github.com/synchthia/altair/logger"
	"github.com/synchthia/altair/nebula"
	"github.com/synchthia/altair/stream"
	"github.com/synchthia/altair/systera"
)

func main() {
	// Init Logger
	logger.Init()

	// Init
	logrus.Printf("[ALTAIR] Starting ALTAIR Bot...")

	// Redis
	redisAddr := os.Getenv("REDIS_ADDRESS")
	if len(redisAddr) == 0 {
		redisAddr = "localhost:6379"
	}
	pool := stream.NewRedisPool(redisAddr)

	// Subscribe
	go func() {
		for {
			stream.PunishmentSubs(pool)
			time.Sleep(3 * time.Second)
		}
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

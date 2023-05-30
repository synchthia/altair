package bot

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var session *discordgo.Session

// InitDiscordBot - Initialize Discord Bot
func InitDiscordBot(token string) error {
	dg, err := discordgo.New("Bot " + token)
	session = dg
	defer dg.Close()

	if err != nil {
		logrus.WithError(err).Errorf("Failed Establish Discord session")
		return err
	}

	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		logrus.WithError(err).Errorf("Failed Open Connection")
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	logrus.Println("[DISCORD] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logrus.Println("[DISCORD] Quitting...")
	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.State.Channel(m.ChannelID)
	guilds, _ := s.State.Guild(ch.GuildID)

	if guilds == nil {
		logrus.Warnf("[Chat] %s tried execute command from out of side: %s", m.Author.Username, m.Content)
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	// Notify Log
	logrus.WithFields(logrus.Fields{
		"username": m.Author.Username,
		"content":  m.Content,
	}).Printf("[Chat]")
}

func hasRole(s *discordgo.Session, m *discordgo.MessageCreate, roleName []string) bool {
	ch, _ := s.State.Channel(m.ChannelID)
	guilds, _ := s.State.Guild(ch.GuildID)

	// if m.Author.ID == guilds.OwnerID {
	// 	return true
	// }

	member, _ := s.State.Member(guilds.ID, m.Author.ID)
	for _, roleID := range member.Roles {
		role, _ := s.State.Role(ch.GuildID, roleID)
		for _, r := range roleName {
			if strings.ToLower(role.Name) == strings.ToLower(r) {
				return true
			}
		}
	}

	return false
}

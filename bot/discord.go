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

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	if strings.Fields(m.Content)[0] == "!servers" {
		if !hasRole(s, m, []string{"administrator"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdServers(s, m)
	}

	if strings.Fields(m.Content)[0] == "!addserver" {
		if !hasRole(s, m, []string{"administrator"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdAddServer(s, m)
	}

	if strings.Fields(m.Content)[0] == "!removeserver" {
		if !hasRole(s, m, []string{"administrator"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdRemoveServer(s, m)
	}

	if strings.Fields(m.Content)[0] == "!lookup" {
		if !hasRole(s, m, []string{"administrator", "developer", "moderator"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdLookup(s, m)
	}

	if strings.Fields(m.Content)[0] == "!announce" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdAnnounce(s, m)
	}

	if strings.Fields(m.Content)[0] == "!dispatch" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdDispatch(s, m)
	}

	if strings.Fields(m.Content)[0] == "!setgroup" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdSetGroup(s, m)
	}

	if strings.Fields(m.Content)[0] == "!creategroup" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdCreateGroup(s, m)
	}

	if strings.Fields(m.Content)[0] == "!removegroup" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdRemoveGroup(s, m)
	}

	if strings.Fields(m.Content)[0] == "!addpermission" || strings.Fields(m.Content)[0] == "!addperms" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdAddPermission(s, m)
	}

	if strings.Fields(m.Content)[0] == "!removepermission" || strings.Fields(m.Content)[0] == "!removeperms" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdRemovePermission(s, m)
	}

	if strings.Fields(m.Content)[0] == "!setmotd" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdSetMotd(s, m)
	}

	if strings.Fields(m.Content)[0] == "!setfavicon" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdSetFavicon(s, m)
	}

	if strings.Fields(m.Content)[0] == "!seen" {
		if !hasRole(s, m, []string{"administrator", "developer"}) {
			s.ChannelMessageSend(m.ChannelID, ":warning: You don't have Permission to do this.")
			return
		}
		cmdSeen(s, m)
	}

	if strings.Fields(m.Content)[0] == "!uuid" {
		cmdUUID(s, m)
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
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

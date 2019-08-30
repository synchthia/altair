package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/synchthia/altair/nebula"
)

func cmdAddServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 4 {
		s.ChannelMessageSend(m.ChannelID, "!addserver <Name> <DisplayName> <Address:Port> [fallback(true/false)]")
		return
	}

	name := args[1]
	displayName := args[2]
	address := args[3]
	fallback := false
	var port int32

	if strings.Contains(args[3], ":") {
		address = strings.Split(args[3], ":")[0]
		i, _ := strconv.ParseInt(strings.Split(args[3], ":")[1], 10, 32)
		port = int32(i)
	} else {
		port = 25565
	}

	if len(args) > 4 {
		if args[4] == "true" {
			fallback = true
		}
	}

	err := nebula.AddServerEntry(name, displayName, address, port, fallback)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Error: `%s`", err.Error()))
		logrus.WithError(err).Errorf("[cmdAddServer]")
		return
	}

	message := fmt.Sprintf(":wrench: **%s (%s)**", displayName, name)
	if fallback {
		message += " :door:"
	}
	embed := &discordgo.MessageEmbed{
		Title: "Server Added: ",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  message,
				Value: fmt.Sprintf("`%s:%d`", address, port),
			},
		},
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdAddServer]")
	}
}

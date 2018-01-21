package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/altair/nebula"
)

func cmdAddServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 4 {
		s.ChannelMessageSend(m.ChannelID, "!addserver <Name> <DisplayName> <Address:Port>")
		return
	}

	name := args[1]
	displayName := args[2]
	address := args[3]
	var port int32

	if strings.Contains(args[3], ":") {
		address = strings.Split(args[3], ":")[0]
		i, _ := strconv.ParseInt(strings.Split(args[3], ":")[1], 10, 32)
		port = int32(i)
	} else {
		port = 25565
	}

	err := nebula.AddServerEntry(name, displayName, address, port)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Error!: %s", err))
		logrus.WithError(err).Errorf("[cmdAddServer]")
	}

	embed := &discordgo.MessageEmbed{
		Title: "Server Added: ",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: fmt.Sprintf(":wrench: **%s (%s)**",
					displayName,
					name,
				),
				Value: fmt.Sprintf("`%s:%d`", address, port),
			},
		},
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdAddServer]")
	}
}

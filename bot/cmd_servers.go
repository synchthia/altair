package bot

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/altair/nebula"
)

func cmdServers(s *discordgo.Session, m *discordgo.MessageCreate) {
	servers, err := nebula.GetServerEntry()
	if err != nil {
		logrus.WithError(err).Errorf("[cmdServers]")
	}

	if len(servers) == 0 {
		s.ChannelMessageSend(m.ChannelID, ":warning: Serverlist is Empty.")
		return
	}

	var fields []*discordgo.MessageEmbedField
	for _, server := range servers {
		field := &discordgo.MessageEmbedField{}
		if server.Status.Online {
			field = &discordgo.MessageEmbedField{
				Name: fmt.Sprintf(":white_check_mark: **%s (%s)** - %d/%d",
					server.DisplayName,
					server.Name,
					server.Status.Players.Online,
					server.Status.Players.Max,
				),
				Value: fmt.Sprintf("`%s:%d`", server.Address, server.Port),
			}
		} else {
			field = &discordgo.MessageEmbedField{
				Name: fmt.Sprintf(":warning: **%s (%s)** - **OFFLINE**",
					server.DisplayName,
					server.Name,
					// server.Status.Players.Online,
					// server.Status.Players.Max,
				),
				Value: fmt.Sprintf("`%s:%d`", server.Address, server.Port),
			}
		}
		fields = append(fields, field)
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Fields: fields})
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdServers]")
	}
}

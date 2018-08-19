package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/altair/nebula"
)

func cmdRemoveServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "!removeserver <Name>")
		return
	}

	name := args[1]

	err := nebula.RemoveServerEntry(name)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Error: `%s`", err.Error()))
		logrus.WithError(err).Errorf("[cmdRemoveServer]")
		return
	}

	_, sendErr := s.ChannelMessageSend(m.ChannelID, "Server Removed: "+name)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdRemoveServer]")
	}
}

package bot

import (
	"fmt"
	"strings"

	"github.com/synchthia/altair/systera"

	"github.com/bwmarrin/discordgo"
)

func cmdUUID(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!uuid <name/uuid>")
		return
	}

	entry, err := systera.GetProfile(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Get %s's Profile: %s", args[1], err.Error()))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Name: %s (UUID: %s)", entry.PlayerName, entry.PlayerUUID))

}

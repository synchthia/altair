package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/altair/systera"
)

func cmdRemoveGroup(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "!removegroup <groupName>")
		return
	}

	err := systera.RemoveGroup(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Remove Group %s : %s", args[1], err.Error()))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Removed Group : %s", args[1]))
}

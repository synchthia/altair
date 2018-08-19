package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/altair/systera"
)

func cmdCreateGroup(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!creategroup <groupName> <prefix>")
		return
	}

	err := systera.CreateGroup(args[1], args[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Create Group %s : %s", args[1], err.Error()))
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Create Group : %s [Prefix: %s]", args[1], args[2]),
	)
}

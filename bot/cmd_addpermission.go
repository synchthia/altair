package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/synchthia/altair/systera"
)

func cmdAddPermission(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 4 {
		s.ChannelMessageSend(m.ChannelID, "!addpermission <groupName> <global/target> <permissions...>")
		return
	}

	err := systera.AddPermission(args[1], args[2], strings.Split(args[3], ","))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Add Permission %s : %s", args[1], err.Error()))
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Added Permissions : %s -> %s: %s", args[1], args[2], strings.Split(args[3], ",")),
	)
}

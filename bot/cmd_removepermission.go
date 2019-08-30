package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/synchthia/altair/systera"
)

func cmdRemovePermission(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "!removepermission <groupName> <global/target> <permissions...>")
		return
	}

	err := systera.RemovePermission(args[1], args[2], strings.Split(args[3], ","))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Remove Permission %s : %s", args[1], err.Error()))
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Remove Permissions : %s -> %s: %s", args[1], args[2], strings.Split(args[3], ",")),
	)
}

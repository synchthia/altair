package bot

import (
	"fmt"
	"strings"

	"github.com/synchthia/altair/nebula"

	"github.com/bwmarrin/discordgo"
)

func cmdSetMotd(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!setmotd <text>")
		return
	}

	message := strings.TrimPrefix(m.Content, args[0]+" ")
	err := nebula.SetMotd(message)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Set Motd: %s", err.Error()))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set Motd: %s", message))
}

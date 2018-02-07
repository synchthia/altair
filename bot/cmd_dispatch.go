package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/altair/systera"
)

func cmdDispatch(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!dispatch <global/[ServerName]> <command>")
		return
	}

	command := strings.TrimPrefix(m.Content, fmt.Sprintf("%s %s ", args[0], args[1]))
	systera.Dispatch(args[1], command)
}

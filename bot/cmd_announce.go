package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/altair/systera"
)

func cmdAnnounce(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!announce <global/[ServerName]> <message>")
		return
	}

	message := strings.TrimPrefix(m.Content, fmt.Sprintf("%s %s ", args[0], args[1]))
	systera.Announce(args[1], message)
}

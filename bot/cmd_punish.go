package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdPunish(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!lookup <name/uuid>")
		return
	}

}

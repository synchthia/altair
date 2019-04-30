package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func cmdShutdown(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@here %s executed Shutdown.", m.Author.Username))
	os.Exit(0)
}

package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/altair/systera"
)

func cmdSetGroup(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	uuid := ""
	groups := []string{}

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!setgroup <name/uuid> [groupA,groupB...]")
		return
	}

	profile, err := systera.GetProfile(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Get %s's Profile: %s", args[1], err.Error()))
		return
	}
	uuid = profile.PlayerUUID

	if len(args) == 2 {
		// groups = append(groups, "default")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's Current Groups %s", profile.PlayerName, profile.Groups))
	} else {
		groups = strings.Split(args[2], ",")
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set %s's Group %s -> %s", profile.PlayerName, profile.Groups, groups))
	systera.SetGroup(uuid, groups)
}

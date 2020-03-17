package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/synchthia/altair/nebula"
)

func cmdLockdown(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "!lockdown <Name> <true|false> [description]")
		return
	}

	name := args[1]
	enabled, parseErr := strconv.ParseBool(args[2])
	if parseErr != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Error: `%s`", parseErr.Error()))
		return
	}

	description := ""
	if enabled && len(args) >= 4 {
		description = strings.TrimPrefix(m.Content, fmt.Sprintf("%s %s %s ", args[0], args[1], args[2]))
		description = strings.ReplaceAll(description, "\\n", "\n")
	}

	e, err := nebula.SetLockdown(name, enabled, description)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Error: `%s`", err.Error()))
		logrus.WithError(err).Errorf("[cmdLockdown]")
		return
	}

	var embed *discordgo.MessageEmbed

	if enabled {
		embed = &discordgo.MessageEmbed{
			Color: 0xE67E22,
			Title: fmt.Sprintf(":lock: %s (%s) is Locked:", e.DisplayName, e.Name),
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  fmt.Sprintf("**%s (%s)**", e.DisplayName, e.Name),
					Value: fmt.Sprintf("```md\n%s\n```", e.Lockdown.Description),
				},
			},
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Color: 0xE67E22,
			Title: fmt.Sprintf(":unlock: %s (%s) is Unlocked", e.DisplayName, e.Name),
		}
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdLockdown]")
	}
}

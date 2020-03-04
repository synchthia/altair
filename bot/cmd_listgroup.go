package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/synchthia/altair/systera"
)

func cmdListGroup(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "!listgroup")
		return
	}

	e, err := systera.GetGroups()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Fetch Groups: %s", err.Error()))
		return
	}

	var fields []*discordgo.MessageEmbedField

	for _, entry := range e.Groups {
		perms := ""
		for _, v := range entry.GlobalPerms {
			perms += fmt.Sprintf("* %s\n", v)
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("__**%s**__", entry.GroupName),
			Value: fmt.Sprintf("** ```md\n# Prefix: %s\n%s\n``` **", entry.GroupPrefix, perms),
		})
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("List of Groups"),
		},
		Fields: fields,
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdListGroup]")
	}

}

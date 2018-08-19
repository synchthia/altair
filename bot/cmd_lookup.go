package bot

import (
	"fmt"
	"strings"

	"gitlab.com/Startail/Systera-API/systerapb"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/altair/systera"
)

func cmdLookup(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!lookup <name/uuid>")
		return
	}

	var r systerapb.PlayerEntry
	r, err := systera.GetProfile(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Get %s's Profile: %s", args[1], err.Error()))
		return
	}

	entries, err := systera.LookupPunish(r.PlayerUUID)
	if err != nil {
		logrus.WithError(err).Errorf("[cmdLookup]")
	}

	if len(entries) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: %s does not have Punish", r.PlayerName))
		return
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Lookup result of %s (%s)", r.PlayerName, r.PlayerUUID),
			IconURL: "https://avatar.minecraft.jp/" + r.PlayerUUID + "/minecraft/m.png",
		},
		Fields: punishFields(entries),
	}
	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdLookup]")
	}
}

package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Systera-API/systerapb"
	"gitlab.com/Startail/altair/systera"
)

// Dummy Entry
type AltLookupDummy struct {
	UUID    string
	Name    string
	Entries []AltLookupEntry
}

type AltLookupEntry struct {
	UUID      string
	Name      string
	IPAddress string
	Hostname  string
	Date      int64
}

func cmdAltLookup(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!altlookup <name/uuid>")
		return
	}

	var r systerapb.PlayerEntry
	r, err := systera.GetProfile(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Get %s's Profile: %s", args[1], err.Error()))
		return
	}

	entries, err := systera.AltLookup(r.PlayerUUID)
	if err != nil {
		logrus.WithError(err).Errorf("[cmdAltLookup]")
	}
	var fields []*discordgo.MessageEmbedField

	for _, entry := range entries {
		addresses := ""
		for i, v := range entry.Addresses {
			addresses += fmt.Sprintf(
				"# %s (%s)\n* FIRST: %s\n   LAST: %s\n",
				v.Hostname,
				v.Address,
				time.Unix(v.FirstSeen/1000, 0).Format("2006-01-02 15:04:05"),
				time.Unix(v.LastSeen/1000, 0).Format("2006-01-02 15:04:05"),
			)
			if i < len(entry.Addresses)-1 {
				addresses += "\t　\n"
			}
		}
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("__**%s**__ - %s", entry.Name, entry.Uuid),
				Value: fmt.Sprintf("** ```md\n%s\n``` **", addresses),
			},
		)
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("AltLookup result of %s (%s)", r.PlayerName, r.PlayerUUID),
			IconURL: "https://avatar.minecraft.jp/" + r.PlayerUUID + "/minecraft/m.png",
		},
		Fields: fields,
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdAltLookup]")
	}

}

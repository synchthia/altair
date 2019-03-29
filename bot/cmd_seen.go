package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/Systera-API/systerapb"
	"gitlab.com/Startail/altair/systera"
)

func cmdSeen(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!seen <name/uuid>")
		return
	}

	var r systerapb.PlayerEntry
	r, err := systera.GetProfile(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Get %s's Profile: %s", args[1], err.Error()))
		return
	}

	// s.ChannelMessageSend(m.ChannelID, ":: "+time.Unix(r.Stats.LastLogin/1000, 0).Format("2006-01-02 15:04:05"))

	currentServer := r.Stats.CurrentServer
	if currentServer == "" {
		currentServer = "Offline"
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Information of %s", r.PlayerName),
			IconURL: "https://avatar.minecraft.jp/" + r.PlayerUUID + "/minecraft/m.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "UUID",
				Value:  fmt.Sprintf("`%s`", r.PlayerUUID),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server",
				Value:  currentServer,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "First Login",
				Value:  time.Unix(r.Stats.FirstLogin/1000, 0).Format("2006-01-02 15:04:05"),
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Last Login",
				Value:  time.Unix(r.Stats.LastLogin/1000, 0).Format("2006-01-02 15:04:05"),
				Inline: false,
			},
		},
	}
	session.ChannelMessageSendEmbed(m.ChannelID, embed)
}

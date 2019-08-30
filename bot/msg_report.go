package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/synchthia/systera-api/systerapb"
)

func ReportMessage(data systerapb.PunishmentStream) {
	roomID := os.Getenv("DISCORD_MODERATOR_ALERT_ROOMID")
	if len(roomID) == 0 {
		return
	}

	entry := data.ReportEntry
	// -> DISCORD_MODERATOR_ALERT_ROOMID
	embed := &discordgo.MessageEmbed{
		Color: 0xFF9800,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Reported (%s -> %s)", entry.From.Name, entry.To.Name),
			IconURL: "https://avatar.minecraft.jp/" + entry.From.Name + "/minecraft/m.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://avatar.minecraft.jp/" + entry.To.Name + "/minecraft/m.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "From",
				Value:  fmt.Sprintf("[%s](%s)", entry.From.Name, "https://minecraft.jp/players/"+entry.From.UUID),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "To",
				Value:  fmt.Sprintf("[%s](%s)", entry.To.Name, "https://minecraft.jp/players/"+entry.To.UUID),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Message",
				Value:  entry.Message,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("- %s | %s", entry.Server, time.Unix(entry.Date/1000, 0).String()),
		},
	}
	session.ChannelMessageSendEmbed(roomID, embed)
}

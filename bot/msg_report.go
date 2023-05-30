package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
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
			IconURL: "https://minotar.net/helm/" + entry.From.Name + "/96.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://minotar.net/helm/" + entry.To.Name + "/96.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Server",
				Value:  entry.Server,
				Inline: false,
			},

			{
				Name:   "From",
				Value:  fmt.Sprintf("[%s](%s)", entry.From.Name, "https://ja.namemc.com/profile/"+entry.From.Uuid),
				Inline: true,
			},
			{
				Name:   "To",
				Value:  fmt.Sprintf("[%s](%s)", entry.To.Name, "https://ja.namemc.com/profile/"+entry.To.Uuid),
				Inline: true,
			},
			{
				Name:   "Message",
				Value:  fmt.Sprintf("```\n%s\n```", entry.Message),
				Inline: false,
			},
		},
		Timestamp: time.UnixMilli(entry.Date).Format(time.RFC3339),
	}
	_, err := session.ChannelMessageSendEmbed(roomID, embed)
	if err != nil {
		logrus.WithError(err).Errorf("Failed parse ReportEntry")
	}
}

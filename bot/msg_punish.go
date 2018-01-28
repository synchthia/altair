package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/Systera-API/systerapb"
)

func PunishMessage(data systerapb.PunishEntryStream) {
	roomID := os.Getenv("DISCORD_MODERATOR_ALERT_ROOMID")
	if len(roomID) == 0 {
		return
	}

	entry := data.Entry
	// -> DISCORD_MODERATOR_ALERT_ROOMID
	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Punished (%s -> %s)", entry.PunishedFrom.Name, entry.PunishedTo.Name),
			IconURL: "https://avatar.minecraft.jp/" + entry.PunishedFrom.Name + "/minecraft/m.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://avatar.minecraft.jp/" + entry.PunishedTo.Name + "/minecraft/m.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "From",
				Value:  fmt.Sprintf("[%s](%s)", entry.PunishedFrom.Name, "https://minecraft.jp/players/"+entry.PunishedFrom.UUID),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "To",
				Value:  fmt.Sprintf("[%s](%s)", entry.PunishedTo.Name, "https://minecraft.jp/players/"+entry.PunishedTo.UUID),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Reason",
				Value:  entry.Reason,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Level",
				Value:  entry.Level.String(),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("- %s | %s", entry.Level, parseExpire(entry.Date, entry.Expire)),
		},
	}
	session.ChannelMessageSendEmbed(roomID, embed)
}

func parseExpire(issued, expire int64) string {
	var msg string
	msg = fmt.Sprintf("%s", time.Unix(issued/1000, 0).String())
	if expire != 0 {
		msg += " / (Expire: " + time.Unix(expire/1000, 0).String() + ")"
	}
	return msg
}

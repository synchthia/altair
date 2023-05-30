package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/synchthia/systera-api/systerapb"
)

func PunishMessage(data systerapb.PunishmentStream) {
	roomID := os.Getenv("DISCORD_MODERATOR_ALERT_ROOMID")
	if len(roomID) == 0 {
		return
	}

	entry := data.PunishStreamEntry.Entry
	// -> DISCORD_MODERATOR_ALERT_ROOMID
	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Punished to %s (%s)", entry.PunishedTo.Name, entry.PunishedTo.Uuid),
			IconURL: "https://minotar.net/helm/" + entry.PunishedFrom.Uuid + "/32.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://minotar.net/helm/" + entry.PunishedTo.Uuid + "/96.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "From",
				Value:  fmt.Sprintf("[%s](%s)", entry.PunishedFrom.Name, "https://ja.namemc.com/profile/"+entry.PunishedFrom.Uuid),
				Inline: true,
			},
			{
				Name:   "To",
				Value:  fmt.Sprintf("[%s](%s)", entry.PunishedTo.Name, "https://ja.namemc.com/profile/"+entry.PunishedTo.Uuid),
				Inline: true,
			},
			{
				Name:  "Level",
				Value: fmt.Sprintf("%s `%s` %s", punishLevelEmoji(entry.Level), entry.Level, punishLevelEmoji(entry.Level)),
			},
			{
				Name: "Expire",
				Value: func() string {
					if (entry.Level == systerapb.PunishLevel_TEMPBAN ||
						entry.Level == systerapb.PunishLevel_PERMBAN) &&
						entry.Expire == 0 {
						return "** FOREVER **"
					} else if entry.Level == systerapb.PunishLevel_WARN ||
						entry.Level == systerapb.PunishLevel_KICK {
						return "-"
					} else {
						return time.UnixMilli(entry.Expire).Format(time.RFC3339)
					}
				}(),
			},
			{
				Name:  "Reason",
				Value: fmt.Sprintf("```\n%s\n```", entry.Reason),
			},
		},
		Timestamp: time.UnixMilli(entry.Date).Format(time.RFC3339),
	}
	session.ChannelMessageSendEmbed(roomID, embed)
}

func punishFields(entries []*systerapb.PunishEntry) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField

	for cnt, entry := range entries {
		expire := "*(Expire: -)*"
		if entry.Expire != 0 {
			expire = fmt.Sprintf("*(Expire: %s)*", time.UnixMilli(entry.Expire).Format("2006-01-02 15:04:05"))
		}
		if entry.Level == systerapb.PunishLevel_PERMBAN {
			expire = "*(Expire: FOREVER)*"
		}

		date := time.UnixMilli(entry.Date).Format("2006-01-02 15:04:05")
		splitter := ""
		if cnt < len(entries)-1 {
			splitter = "\n⠀"
		}

		field := &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s **%s** : `%s`", punishLevelEmoji(entry.Level), entry.Level, entry.Reason),
			Value: fmt.Sprintf("- %s by %s\n %s\n%s", date, entry.PunishedFrom.Name, expire, splitter),
		}

		fields = append(fields, field)
	}

	return fields
}

func punishLevelEmoji(level systerapb.PunishLevel) string {
	switch level {
	case systerapb.PunishLevel_WARN:
		return ":loudspeaker:"
	case systerapb.PunishLevel_KICK:
		return ":raised_hand:"
	case systerapb.PunishLevel_TEMPBAN:
		return ":octagonal_sign:"
	case systerapb.PunishLevel_PERMBAN:
		return ":no_entry:"
	default:
		return ":question:"
	}
}

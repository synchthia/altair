package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Startail/Systera-API/systerapb"
)

func PunishMessage(data systerapb.PunishmentStream) {
	roomID := os.Getenv("DISCORD_MODERATOR_ALERT_ROOMID")
	if len(roomID) == 0 {
		return
	}

	entry := data.PunishEntry
	entries := []*systerapb.PunishEntry{entry}
	// -> DISCORD_MODERATOR_ALERT_ROOMID
	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Punished %s (%s)", entry.PunishedTo.Name, entry.PunishedTo.UUID),
			IconURL: "https://avatar.minecraft.jp/" + entry.PunishedTo.UUID + "/minecraft/s.png",
		},
		Fields: punishFields(entries),
	}
	session.ChannelMessageSendEmbed(roomID, embed)
}

func punishFields(entries []*systerapb.PunishEntry) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField

	for cnt, entry := range entries {
		expire := "*(Expire: -)*"
		if entry.Expire != 0 {
			expire = fmt.Sprintf(" *(Expire: %s)*", time.Unix(entry.Expire/1000, 0).Format("2006-01-02 15:04:05"))
		}
		if entry.Level == systerapb.PunishLevel_PERMBAN {
			expire = "*(Expire: FOREVER)*"
		}

		date := time.Unix(entry.Date/1000, 0).Format("2006-01-02 15:04:05")
		splitter := ""
		if cnt < len(entries)-1 {
			splitter = "\n⠀"
		}

		field := &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s **%s** ≫ `%s`", punishLevelEmoji(entry.Level), entry.Level, entry.Reason),
			Value: fmt.Sprintf("- %s by %s\n%s\n%s", date, entry.PunishedFrom.Name, expire, splitter),
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

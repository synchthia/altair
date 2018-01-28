package bot

import (
	"fmt"
	"strings"
	"time"

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
	r, err := getProfile(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Lookup Player's Profile: %s", err.Error()))
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
			splitter = "<:blank:407214259815186433>"
		}

		field := &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s **%s** (by %s) - %s", parseLevelEmoji(entry.Level), entry.Level, entry.PunishedFrom.Name, date),
			Value: fmt.Sprintf("- %s\n- `%s`\n%s", expire, entry.Reason, splitter),
		}
		//msg := fmt.Sprintf("**%s** > [%s]%s [(by %s)]%s\n", date, entry.Level, entry.Reason, entry.PunishedFrom.Name, expire)
		fields = append(fields, field)
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x009688,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Lookup result of %s (%s)", r.PlayerName, r.PlayerUUID),
			IconURL: "https://avatar.minecraft.jp/" + r.PlayerUUID + "/minecraft/m.png",
		},
		Fields: fields,
	}
	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if sendErr != nil {
		logrus.WithError(sendErr).Errorf("[cmdLookup]")
	}
}

func getProfile(uuidOrString string) (systerapb.PlayerEntry, error) {
	if len(uuidOrString) == 32 {
		r, err := systera.FetchPlayerProfile(uuidOrString)
		if err != nil {
			logrus.WithError(err).Errorf("[cmdLookup] Failed Lookup Player's Profile(from UUID)")
			return systerapb.PlayerEntry{}, err
		}
		return *r, nil
	}
	r, err := systera.FetchPlayerProfileByName(uuidOrString)
	if err != nil {
		logrus.WithError(err).Errorf("[cmdLookup] Failed Lookup Player's Profile(from Name)")
		return systerapb.PlayerEntry{}, err
	}
	return *r, nil
}

func parseLevelEmoji(level systerapb.PunishLevel) string {
	switch level {
	case systerapb.PunishLevel_WARN:
		return ":thinking:"
	case systerapb.PunishLevel_KICK:
		return ":angry:"
	case systerapb.PunishLevel_TEMPBAN:
		return ":rage:"
	case systerapb.PunishLevel_PERMBAN:
		return ":innocent:"
	default:
		return ":question:"
	}
}

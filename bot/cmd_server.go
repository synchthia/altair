package bot

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/synchthia/altair/nebula"
	"github.com/synchthia/nebula-api/nebulapb"
)

func cmdServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "!server <Name>")
		return
	}
	name := args[1]

	servers, err := nebula.GetServerEntry()
	if err != nil {
		logrus.WithError(err).Errorf("[cmdServer]")
	}

	if len(servers) == 0 {
		s.ChannelMessageSend(m.ChannelID, ":warning: Serverlist is Empty.")
		return
	}

	found := false
	for _, server := range servers {
		if server.Name != name {
			continue
		}

		var entries []*discordgo.MessageEmbedField

		entries = append(entries, &discordgo.MessageEmbedField{
			Name:  "**Name**",
			Value: server.Name,
		}, &discordgo.MessageEmbedField{
			Name:  "**DisplayName**",
			Value: server.DisplayName,
		}, &discordgo.MessageEmbedField{
			Name:  "**Address**",
			Value: fmt.Sprintf("`%s:%d`", server.Address, server.Port),
		})

		if server.Status.Online {
			entries = append(entries, &discordgo.MessageEmbedField{
				Name:  "**Players**",
				Value: fmt.Sprintf("%d/%d", server.Status.Players.Online, server.Status.Players.Max),
			}, &discordgo.MessageEmbedField{
				Name:  "**Description**",
				Value: fmtDescription(server.Status.Description),
			}, &discordgo.MessageEmbedField{
				Name:  "**Version**",
				Value: fmt.Sprintf("%s (%d)", server.Status.Version.Name, server.Status.Version.Protocol),
			})

		}

		entries = append(entries, &discordgo.MessageEmbedField{
			Name:   "**Status**",
			Value:  fmtStatus(server.Status.Online),
			Inline: true,
		}, &discordgo.MessageEmbedField{
			Name:   "**Fallback**",
			Value:  fmtFallback(server.Fallback),
			Inline: true,
		}, &discordgo.MessageEmbedField{
			Name:   "**Lockdown**",
			Value:  fmtLockdown(server.Lockdown),
			Inline: true,
		})

		embed := &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("%s (%s)", server.DisplayName, server.Name),
			Fields: entries,
		}

		_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if sendErr != nil {
			logrus.WithError(sendErr).Errorf("[cmdServer]")
		}
		found = true
	}

	if !found {
		s.ChannelMessageSend(m.ChannelID, ":warning: Server not found: "+name)
	}
}

func fmtDescription(description string) string {
	if len(description) == 0 {
		return "-"
	}

	return fmt.Sprintf("`%s`", description)
}

func fmtFallback(fallback bool) string {
	if fallback {
		return ":door: ON"
	}
	return ":no_entry_sign: OFF"
}

func fmtLockdown(l *nebulapb.Lockdown) string {
	if l.Enabled {
		return fmt.Sprintf(":lock: Locked:\n `%s`", l.Description)
	}
	return ":unlock: Not Locked"
}

func fmtStatus(online bool) string {
	if online {
		return ":white_check_mark: Online"
	}
	return ":warning: Offline"
}

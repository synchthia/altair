package bot

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/synchthia/altair/nebula"

	"github.com/bwmarrin/discordgo"
)

func cmdSetFavicon(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "!setfavicon <url>")
		return
	}

	url := args[1]
	resp, err1 := http.Get(url)
	if err1 != nil {
		logrus.WithError(err1).Errorf("[SetFavicon]")
		return
	}

	defer resp.Body.Close()

	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		logrus.WithError(err2).Errorf("[SetFavicon]")
		return
	}

	encoded := base64.StdEncoding.EncodeToString(b)

	fmt.Print()
	err := nebula.SetFavicon(fmt.Sprintf("data:image/png;base64,%s", encoded))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":warning: Failed Set Favicon: %s", err.Error()))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set Favicon from: %s", url))
}

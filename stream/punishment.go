package stream

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/synchthia/altair/bot"
	"github.com/synchthia/systera-api/systerapb"
)

// PunishmentSubs - Subscribe Report Stream
func PunishmentSubs(pool *redis.Pool) error {
	c := pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	err := psc.PSubscribe("systera.punishment.global")
	if err != nil {
		logrus.WithError(err).Errorf("[REDIS] Error while Subscribe Punishment")
		return err
	}

	for {
		switch n := psc.Receive().(type) {
		case redis.PMessage:
			logrus.Printf("[PUNISH] Incoming: %s %s", n.Pattern, n.Channel)
			entry := systerapb.PunishmentStream{}
			json.Unmarshal(n.Data, &entry)

			switch entry.Type {
			case systerapb.PunishmentStream_PUNISH:
				bot.PunishMessage(entry)
				break
			case systerapb.PunishmentStream_REPORT:
				bot.ReportMessage(entry)
				break
			}

		case redis.Subscription:
			logrus.Printf("[PUNISH] Subscription: %s %s %d", n.Kind, n.Channel, n.Count)
			if n.Count == 0 {
				return nil
			}
		case error:
			logrus.WithError(n).Errorf("[PUNISH] Error")
			return err
		}
	}
}

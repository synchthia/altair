package stream

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/synchthia/systera-api/systerapb"
	"github.com/synchthia/altair/bot"
)

// PunishmentSubs - Subscribe Report Stream
func PunishmentSubs() {
	c := pool.Get()
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	psc := redis.PubSubConn{Conn: c}

	go func() {
		defer wg.Done()
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
					return
				}
			case error:
				fmt.Printf("error: %v\n", n)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		err := psc.PSubscribe("systera.punishment.global")
		if err != nil {
			logrus.WithError(err).Errorf("[REDIS] Error while Subscribe Punishment")
			return
		}
	}()
	wg.Wait()
}

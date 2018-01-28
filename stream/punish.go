package stream

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Systera-API/systerapb"
	"gitlab.com/Startail/altair/bot"
)

// PunishSubs - Subscribe Report Stream
func PunishSubs() {
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
				entry := systerapb.PunishEntryStream{}
				json.Unmarshal(n.Data, &entry)

				bot.PunishMessage(entry)

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
		err := psc.PSubscribe("systera.punish")
		if err != nil {
			logrus.WithError(err).Errorf("[REDIS] Error while Subscribe Punish")
			return
		}
	}()
	wg.Wait()
}

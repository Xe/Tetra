package tetra

import (
	"time"
)

// Continuously reports the network statistics. Should be run in a
// gorotutine.
func (t *Tetra) GetNetworkStats() {
	for {
		num := int64(len(t.Clients.ByNick))

		if num == 0 {
			wait()
		}

		t.Clients.Gauge.Update(num)

		for _, server := range t.Servers {
			if server.Counter == nil {
				continue
			}

			server.Counter.Update(int64(server.count))
		}

		t.Log.Printf("Logged stats for network and server populations")

		wait()
	}
}

func (t *Tetra) GetChannelStats() {
	for {
		for _, channel := range t.Channels {
			channel.Gauge.Update(int64(len(channel.Clients)))
		}

		t.Log.Printf("Logged stats for channel populations")

		wait()
	}
}

func wait() {
	time.Sleep(time.Minute * 5)
}

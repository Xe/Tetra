package tetra

// Continuously reports the network statistics. Should be run in a
// gorotutine.
func (t *Tetra) GetNetworkStats(...interface{}) {
	num := int64(len(t.Clients.ByNick))

	if num == 0 {
		return
	}

	t.Clients.Gauge.Update(num)

	for _, server := range t.Servers {
		if server.Counter == nil {
			continue
		}

		server.Counter.Update(int64(server.count))
	}

	debug("Logged stats for network and server populations")
}

func (t *Tetra) GetChannelStats(...interface{}) {
	for _, channel := range t.Channels {
		channel.Gauge.Update(int64(len(channel.Clients)))
	}

	debug("Logged stats for channel populations")
}

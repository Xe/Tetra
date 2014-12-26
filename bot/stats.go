package tetra

func GetNetworkStats(...interface{}) {
	num := int64(len(Clients.ByNick))

	if num == 0 {
		return
	}

	Clients.Gauge.Update(num)

	for _, server := range Servers {
		if server.Counter == nil {
			continue
		}

		server.Counter.Update(int64(server.count))
	}

	debug("Logged stats for network and server populations")
}

func GetChannelStats(...interface{}) {
	for _, channel := range Channels {
		channel.Gauge.Update(int64(len(channel.Clients)))
	}

	debug("Logged stats for channel populations")
}

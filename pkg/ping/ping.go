// ping package implements the mincraft Server List Ping protocol.
// This package is strictly compliant with the following documentation : https://wiki.vg/Server_List_Ping
package ping

// Ping returns the server list ping infos (JSON-like object), and latency of a minecraft server.
// If an error occured at any point of the process, an nil json response, a latency of -1, and a non nil error are returned.
func Ping(hostname string, port int) (map[string]interface{}, int, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return nil, -1, err
	}

	handshake, err := client.Handshake()
	if err != nil {
		return nil, -1, err
	}

	latency, err := client.Ping()
	if err != nil {
		return nil, -1, err
	}

	err = client.Disconnect()
	if err != nil {
		return nil, -1, err
	}

	return handshake.Properties, latency, nil
}

// PingLegacy returns the legacy server list ping infos, and latency of a minecraft server.
// If an error occured at any point of the process, an empty response, a latency of -1, and a non nil error are returned.
func PingLegacy(hostname string, port int) (LegacyPingInfos, int, error) {
	client := NewClientLegacy(hostname, port)

	err := client.Connect()
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	infos, latency, err := client.Ping()
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	err = client.Disconnect()
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	return infos, latency, nil
}

// PingLegacy1_6_4 returns the legacy server list ping infos (with 1.6+ SLP), and latency of a minecraft server.
// If an error occured at any point of the process, an empty response, a latency of -1, and a non nil error are returned.
func PingLegacy1_6_4(hostname string, port int) (LegacyPingInfos, int, error) {
	client := NewClientLegacy(hostname, port)

	err := client.Connect()
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	infos, latency, err := client.Ping1_6_4()
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	err = client.Disconnect()
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	return infos, latency, nil
}

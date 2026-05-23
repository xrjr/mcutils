// ping package implements the mincraft Server List Ping protocol.
// This package is strictly compliant with the following documentation : https://minecraft.wiki/w/Java_Edition_protocol/Server_List_Ping.
package ping

import "errors"

// Ping returns the server list ping infos (JSON-like object), and latency of a minecraft server.
// If an error occurred at any point of the process, an nil json response, a latency of -1, and a non nil error are returned.
// If the server responds to the ping request with a bad packet (e.g. with a handshake response), the packet will not be read and the error will be ingored (to support Forge servers).
func Ping(hostname string, port int) (JSON, int, error) {
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

	// Some forge servers respond to ping request with the handshake response. In this case, a ErrInvalidPacketType will be returned.
	// We'll be ingoring this error because it doesn't have any side effect, since :
	//   - we don't retrieve any information from the pong response packet
	//   - connection is closed right after
	if err != nil && !errors.Is(err, ErrInvalidPacketType) {
		return nil, -1, err
	}

	err = client.Disconnect()
	if err != nil {
		return nil, -1, err
	}

	return handshake.Properties, latency, nil
}

// PingLegacy returns the legacy server list ping infos, and latency of a minecraft server.
// If an error occurred at any point of the process, an empty response, a latency of -1, and a non nil error are returned.
// If the minecraft server has a version <= 1.3, ProtocolNumber and MinecraftVersion are not set.
// Note that legacy ping should be working on most servers that don't require host to be set, but it is notoriously slow on 1.6.x vanilla servers.
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

// PingLegacy1_6_4 returns the legacy server list ping infos (using 1.6+ SLP protocol), and latency of a minecraft server.
// If an error occurred at any point of the process, an empty response, a latency of -1, and a non nil error are returned.
// Note that on vanilla servers < 1.4.x, this protocol usually don't work.
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

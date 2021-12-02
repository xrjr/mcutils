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

// bedrock package implements the unconnected ping sequence of the raknet protocol (used by minecraft bedrock servers).
// This package is strictly compliant with the following documentation : https://wiki.vg/Raknet_Protocol.
package bedrock

// Ping returns the server infos, and latency of a minecraft bedrock server.
// If an error occurred at any point of the process, an empty pong response, a latency of -1, and a non nil error are returned.
func Ping(hostname string, port int) (UnconnectedPong, int, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return UnconnectedPong{}, -1, err
	}

	unconnectedPong, latency, err := client.UnconnectedPing()
	if err != nil {
		return UnconnectedPong{}, -1, err
	}

	err = client.Disconnect()
	if err != nil {
		return UnconnectedPong{}, -1, err
	}

	return unconnectedPong, latency, nil
}

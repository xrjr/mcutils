// query package implements the mincraft query protocol.
// This package is strictly compliant with the following documentation : https://wiki.vg/Query
package query

// QueryBasic returns the basic stat of a minecraft server.
// If an error occured at any point of the process, an empty BasicStat and a non nil error are returned.
func QueryBasic(hostname string, port int) (BasicStat, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return BasicStat{}, err
	}

	token, err := client.Handshake()
	if err != nil {
		return BasicStat{}, err
	}

	basicStat, err := client.BasicStat(token)
	if err != nil {
		return BasicStat{}, err
	}

	err = client.Disconnect()
	if err != nil {
		return BasicStat{}, err
	}

	return basicStat, nil
}

// QueryFull returns the full stat of a minecraft server.
// If an error occured at any point of the process, an empty FullStat and a non nil error are returned.
func QueryFull(hostname string, port int) (FullStat, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return FullStat{}, err
	}

	token, err := client.Handshake()
	if err != nil {
		return FullStat{}, err
	}

	fullStat, err := client.FullStat(token)
	if err != nil {
		return FullStat{}, err
	}

	err = client.Disconnect()
	if err != nil {
		return FullStat{}, err
	}

	return fullStat, nil
}

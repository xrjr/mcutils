package bedrock

// Stat returns the status of the minecraft:bedrock server
func Stat(hostname string, port int) (BEStat, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return BEStat{}, err
	}

	stat, err := client.Stat()
	if err != nil {
		return BEStat{}, err
	}

	err = client.Disconnect()
	if err != nil {
		return BEStat{}, err
	}

	return *stat, nil
}

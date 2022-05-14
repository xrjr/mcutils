package bedrock

import (
	"fmt"
	"os"
)

// Stat returns the status of the minecraft:bedrock server
func Stat(hostname string, port int) (BEStat, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return BEStat{}, err
	}

	defer func() {
		err = client.Disconnect()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Disconnect error: ", err.Error())
		}
	}()

	stat, err := client.Stat()
	if err != nil {
		return BEStat{}, err
	}

	return *stat, nil
}

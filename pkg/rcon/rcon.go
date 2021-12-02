// rcon implements the minecraft rcon protocol, which is itself an implementation of the source rcon protocol.
// This package is strictly compliant with the following documentation : https://wiki.vg/RCON
package rcon

import "errors"

var (
	ErrWrongPassword error = errors.New("wrong password")
)

// Rcon executes a command on a minecraft server, and returns the response of that command.
// If the password is wrong, the error will be of type ErrWrongPassword.
func Rcon(hostname string, port int, password string, command string) (string, error) {
	client := NewClient(hostname, port)

	err := client.Connect()
	if err != nil {
		return "", err
	}

	ok, err := client.Authenticate(password)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrWrongPassword
	}

	response, err := client.Command(command)
	if err != nil {
		return "", nil
	}

	client.Disconnect()
	if err != nil {
		return "", err
	}

	return response, nil
}

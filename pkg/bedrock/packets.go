package bedrock

type BEStat struct {
	GameName        string
	MOTD            string
	HostName        string
	Protocol        string
	Version         string
	Players         int
	MaxPlayers      int
	ServerID        int64
	Map             string
	GameMode        string
	NintendoLimited string
	IPv4Port        int
	IPv6Port        int
	Extra           string
}

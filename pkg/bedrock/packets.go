package bedrock

// unconnectedPongResponse is the type respresenting the response of the unconnected ping request.
type unconnectedPongResponse struct {
	PacketID        byte
	ClientTimestamp uint64
	ServerGUID      uint64
	Magic           []byte

	GameName         string
	MOTD             string
	ProtocolVersion  int
	MinecraftVersion string
	OnlinePlayers    int
	MaxPlayers       int
	ServerID         string
	Map              string
	GameMode         string
	IPv4Port         int
	IPv6Port         int
}

// unconnectedPong transforms the unconnectedPongResponse into a more human-usable UnconnectedPong struct.
func (upr *unconnectedPongResponse) unconnectedPong() UnconnectedPong {
	return UnconnectedPong{
		GameName:         upr.GameName,
		MOTD:             upr.MOTD,
		ProtocolVersion:  upr.ProtocolVersion,
		MinecraftVersion: upr.MinecraftVersion,
		OnlinePlayers:    upr.OnlinePlayers,
		MaxPlayers:       upr.MaxPlayers,
		ServerID:         upr.ServerID,
		Map:              upr.Map,
		GameMode:         upr.GameMode,
		IPv4Port:         upr.IPv4Port,
		IPv6Port:         upr.IPv6Port,
	}
}

// UnconnectedPong contains unconnected pong informations.
type UnconnectedPong struct {
	GameName         string
	MOTD             string
	ProtocolVersion  int
	MinecraftVersion string
	OnlinePlayers    int
	MaxPlayers       int
	ServerID         string
	Map              string
	GameMode         string
	IPv4Port         int
	IPv6Port         int
}

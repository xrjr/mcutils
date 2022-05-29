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
	LevelName        string
	GameMode         string
	GameModeNumeric  int
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
		LevelName:        upr.LevelName,
		GameMode:         upr.GameMode,
		GameModeNumeric:  upr.GameModeNumeric,
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
	LevelName        string
	GameMode         string
	GameModeNumeric  int
	IPv4Port         int
	IPv6Port         int
}

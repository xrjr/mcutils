package query

// packet is the common structure conatained in all query datagrams
type packet struct {
	Type      byte
	SessionID uint32
}

// handshakeResponse is the type respresenting the response of the handshake query
type handshakeResponse struct {
	packet
	ChallengeToken uint32
}

// basicStatResponse is the type respresenting the response of the basic stat query
type basicStatResponse struct {
	packet
	MOTD       string
	GameType   string
	Map        string
	NumPlayers int
	MaxPlayers int
	HostPort   int16
	HostIP     string
}

// basicStat transforms the basicStatResponse into a more human-usable BasicStat struct.
func (bsr *basicStatResponse) basicStat() BasicStat {
	return BasicStat{
		MOTD:       bsr.MOTD,
		GameType:   bsr.GameType,
		Map:        bsr.Map,
		NumPlayers: bsr.NumPlayers,
		MaxPlayers: bsr.MaxPlayers,
		HostPort:   int(bsr.HostPort),
		HostIP:     bsr.HostIP,
	}
}

// fullStatResponse is the type respresenting the response of the full stat query
type fullStatResponse struct {
	packet
	Padding1       [11]byte
	KVSection      map[string]string
	Padding2       [10]byte
	PlayersSection []string
}

// fullStat transforms the fullStatResponse into a more human-usable FullStat struct.
func (fsr *fullStatResponse) fullStat() FullStat {
	return FullStat{
		Properties:    fsr.KVSection,
		OnlinePlayers: fsr.PlayersSection,
	}
}

// BasicStat is the user-friendly/user-returnable version of basicStatResponse
type BasicStat struct {
	MOTD       string `json:"motd"`
	GameType   string `json:"gameType"`
	Map        string `json:"map"`
	NumPlayers int    `json:"numPlayers"`
	MaxPlayers int    `json:"maxPlayers"`
	HostPort   int    `json:"hostPort"`
	HostIP     string `json:"hostIp"`
}

// FullStat is the user-friendly/user-returnable version of fullStatResponse
type FullStat struct {
	Properties    map[string]string `json:"properties"`
	OnlinePlayers []string          `json:"onlinePlayers"`
}

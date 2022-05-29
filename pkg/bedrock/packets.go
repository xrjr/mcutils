package bedrock

type unconnectedPongResponse struct {
	ID              byte
	ClientTimestamp uint64
	ServerGUID      uint64
	Magic           []byte
	Data            string
}

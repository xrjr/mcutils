package bedrock

type unconnectedPongResponse struct {
	PacketID        byte
	ClientTimestamp uint64
	ServerGUID      uint64
	Magic           []byte
	Data            string
}

package peers

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// Peer encodes connection information for a peer
type Peer struct {
	IP     net.IP
	Port   uint16
	Tries  uint32
	PeerID [20]byte
}

// Unmarshal parses peer IP addresses and ports from a buffer
func Unmarshal(peerID [20]byte, peersBin []byte) ([]Peer, error) {
	const peerSize = 6 // 4 for IP, 2 for port
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBin[offset+4 : offset+6]))
		peers[i].Tries = 0
		peers[i].PeerID = peerID
	}
	return peers, nil
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

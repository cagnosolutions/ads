package bpt

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"time"
)

func init() {
	cachedNodeId = getNodeId()
}

var cachedNodeId []byte

func getNodeId() []byte {
	var d [6]byte
	inet, err := net.Interfaces()
	if err == nil {
		var set bool
		for _, v := range inet {
			if len(v.HardwareAddr.String()) != 0 {
				copy(d[:], []byte(v.HardwareAddr))
				set = true
				break
			}
		}
		if set {
			return d[:]
		}
	}
	rand.Read(d[:])
	d[0] |= 0x01
	return d[:]
}

func UUID() []byte {
	t := uint64(time.Now().UnixNano()/100 + 0x01b21dd213814000)
	var b [2]byte
	rand.Read(b[:])
	clockSeq := binary.LittleEndian.Uint16(b[:])
	clockSeq &= 0x3FFF
	u := make([]byte, 16)
	binary.LittleEndian.PutUint32(u[0:4], uint32(t&(0x100000000-1)))
	binary.LittleEndian.PutUint16(u[4:6], uint16((t>>32)&0xFFFF))
	binary.LittleEndian.PutUint16(u[6:8], uint16((t>>48)&0x0FFF))
	binary.LittleEndian.PutUint16(u[8:10], clockSeq)
	copy(u[10:16], cachedNodeId)
	u[8] &= 0x3F
	u[8] |= 0x80
	u[6] = (u[6] & 0x0F) | (0x01 << 4)
	return u
}

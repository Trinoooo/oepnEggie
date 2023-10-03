package archive

import "encoding/binary"

func parsePort(port int) int {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(port))
	return int(binary.BigEndian.Uint16(buf))
}

func parseAddr() [4]byte {
	addr := uint32(0) // 0.0.0.0
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, addr)
	binary.BigEndian.Uint32(buf)
	res := [4]byte{}
	for idx, b := range buf {
		res[idx] = b
	}
	return res
}

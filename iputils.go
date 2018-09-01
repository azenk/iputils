package iputils

import (
	"errors"
	"net"
)

var (
	ErrOutOfRange = errors.New("specified offset and/or width are invalid")
)

// SetBits sets the bits starting at startOffset and ending at startOffset + setWidth
// to the bits specified.
func SetBits(ip net.IP, bits uint64, startOffset, setWidth uint) (net.IP, error) {
	var newIP []byte

	var size uint
	if v4ip := ip.To4(); v4ip != nil {
		ip = v4ip
		newIP = make([]byte, 4)
		copy(newIP, net.IPv4zero)
		size = 32
	} else {
		newIP = make([]byte, 16)
		copy(newIP, net.IPv6zero)
		size = 128
	}

	if startOffset+setWidth > size {
		return newIP, ErrOutOfRange
	}

	maskBits := ^uint64(0) >> uint(64-setWidth)
	bits &= maskBits

	for i, b := range ip {
		newIP[i] = b
		if (i+1)*8 > int(startOffset) && i*8 < int(startOffset+setWidth) {
			shift := int(startOffset+setWidth) - (i+1)*8
			var t, mask byte
			if shift > 0 {
				mask = byte(maskBits >> uint(shift) & 0x0ff)
				t = byte(bits>>uint(shift)) & mask
			} else {
				mask = byte(maskBits << uint(-1*shift) & 0x0ff)
				t = byte(bits<<uint(-1*shift)) & mask
			}
			newIP[i] = b&^mask | t
		}
	}
	return newIP, nil
}

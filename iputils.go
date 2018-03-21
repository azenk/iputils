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
	newIP := make([]byte, 16)

	if ip.To4() != nil {
		copy(newIP, net.IPv4zero)
		startOffset = startOffset + 96
	} else {
		copy(newIP, net.IPv6zero)
	}

	if startOffset+setWidth > 128 {
		return newIP, ErrOutOfRange
	}

	bits &= ^uint64(0) >> uint(64-setWidth)

	for i, b := range ip {
		newIP[i] = b
		if (i+1)*8 > int(startOffset) && i*8 < int(startOffset+setWidth) {
			shift := int(startOffset+setWidth) - (i+1)*8
			var t, mask byte
			if shift > 0 {
				mask = byte(^uint64(0) >> uint(shift) & 0x0ff)
				t = byte(bits>>uint(shift)) & mask
			} else {
				mask = byte(^uint64(0) << uint(-1*shift) & 0x0ff)
				t = byte(bits<<uint(-1*shift)) & mask
			}
			newIP[i] = b&mask | t
		}
	}
	return newIP, nil
}

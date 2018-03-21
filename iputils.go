package iputils

import (
	"net"
)

func SetBits(ip net.IP, bits uint64, startOffset, setWidth int) (net.IP, error) {
	var newIP net.IP

	if ip.To4() != nil {
		newIP = net.IPv4zero
		startOffset = startOffset + 96
	} else {
		newIP = net.IPv6zero
	}

	for i, b := range ip {
		newIP[i] = b
		if (i+1)*8 > startOffset && i*8 <= startOffset+setWidth {
			shift := startOffset + setWidth - (i+1)*8
			var t, mask byte
			if shift > 0 {
				t = byte(bits >> uint(shift) & 0x00ff)
				mask = byte(0x0ff >> uint(shift) & 0x0ff)
			} else {
				t = byte(bits << uint(-1*shift) & 0x00ff)
				mask = byte(0x0ff << uint(-1*shift) & 0x0ff)
			}
			newIP[i] = b&mask | t
		}
	}
	return newIP, nil
}

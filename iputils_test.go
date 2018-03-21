package iputils

import (
	"net"
	"testing"
)

func TestSetBits(t *testing.T) {
	type ipSetTestCase struct {
		BaseIP     net.IP
		Bits       uint64
		Offset     int
		Width      int
		ExpectedIP net.IP
	}

	testCases := []ipSetTestCase{
		ipSetTestCase{net.ParseIP("10.0.0.0"), 0xff, 16, 8, net.ParseIP("10.0.255.0")},
		ipSetTestCase{net.ParseIP("10.0.15.0"), 0xff, 12, 8, net.ParseIP("10.15.240.0")},
		ipSetTestCase{net.ParseIP("2001:db8:0:1::"), 0xffffffff, 64, 32, net.ParseIP("2001:db8:0:1:ffff:ffff::")},
		ipSetTestCase{net.ParseIP("2001:db8:0:1::"), 0xffffffff, 66, 32, net.ParseIP("2001:db8:0:1:3fff:ffff:c000:0")},
		ipSetTestCase{net.ParseIP("2001:db8:0:1::"), 0xf3f2f1f0, 72, 24, net.ParseIP("2001:db8:0:1:00f2:f1f0::")},
	}

	for _, testCase := range testCases {
		newIP, err := SetBits(testCase.BaseIP, testCase.Bits, testCase.Offset, testCase.Width)
		if err != nil {
			t.Fatalf("unable to set bits: %v", err)
		}

		if newIP.String() != testCase.ExpectedIP.String() {
			t.Errorf("Expected IP: %s, Got: %s", testCase.ExpectedIP.String(), newIP.String())
		}
	}
}

package iputils

import (
	"net"
	"testing"
)

func TestSetBits(t *testing.T) {
	type ipSetTestCase struct {
		BaseIP      net.IP
		Bits        uint64
		Offset      uint
		Width       uint
		ExpectedIP  net.IP
		ExpectedErr error
	}

	testCases := []ipSetTestCase{
		ipSetTestCase{net.ParseIP("10.0.0.0"), 0xff, 16, 8, net.ParseIP("10.0.255.0"), nil},
		ipSetTestCase{net.ParseIP("10.2.3.64").To4(), 0x8e28f0ef6bafe9fb, 28, 4, net.ParseIP("10.2.3.75"), nil},
		ipSetTestCase{net.ParseIP("10.0.15.0"), 0xff0, 12, 12, net.ParseIP("10.15.240.0"), nil},
		ipSetTestCase{net.ParseIP("2001:db8:0:1::"), 0xffffffff, 64, 32, net.ParseIP("2001:db8:0:1:ffff:ffff::"), nil},
		ipSetTestCase{net.ParseIP("2001:db8:0:1::"), 0xffffffff, 66, 32, net.ParseIP("2001:db8:0:1:3fff:ffff:c000:0"), nil},
		ipSetTestCase{net.ParseIP("2001:db8:0:1:00ff::"), 0xf3f2f1f0, 72, 24, net.ParseIP("2001:db8:0:1:00f2:f1f0::"), nil},
		ipSetTestCase{net.ParseIP("2001:db8:0:6:683f:ea00::"), ^uint64(0), 92, 36, net.ParseIP("2001:db8:0:6:683f:ea0f:ffff:ffff"), nil},
		ipSetTestCase{net.ParseIP("2001:db8:0:1::"), 0xf3f2f1f0, 100, 32, net.ParseIP("::"), ErrOutOfRange},
		ipSetTestCase{net.ParseIP("10.0.0.0"), 0xff, 25, 8, net.ParseIP("0.0.0.0"), ErrOutOfRange},
	}

	for _, testCase := range testCases {
		newIP, err := SetBits(testCase.BaseIP, testCase.Bits, testCase.Offset, testCase.Width)
		if err != testCase.ExpectedErr {
			t.Errorf("Expected err %v, got: %v", testCase.ExpectedErr, err)
		}

		if newIP.String() != testCase.ExpectedIP.String() {
			t.Errorf("Expected IP: %s, Got: %s", testCase.ExpectedIP.String(), newIP.String())
		}
	}
}

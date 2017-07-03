package random

import (
	"testing"
)

func Test_RandomService_Bit(t *testing.T) {
	var zeroFound bool
	var oneFound bool

	for i := 0; i < 100; i++ {
		b := Bit()

		if b == 0 {
			zeroFound = true
		}
		if b == 1 {
			oneFound = true
		}
		if zeroFound && oneFound {
			break
		}
	}

	if !zeroFound || !oneFound {
		t.Fatal("expected", false, "got", true)
	}
}

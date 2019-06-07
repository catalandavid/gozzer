package misc

import (
	"testing"
	"time"
)

func forceReturnTrue() bool {
	return true
}

func forceReturnFalse() bool {
	return false
}

func TestPollCheckUntilForcePass(t *testing.T) {
	res := PollCheckUntil(forceReturnFalse, 10*time.Millisecond)
	if res == nil {
		t.Fail()
	}
}

func TestPollCheckUntilForceFail(t *testing.T) {
	res := PollCheckUntil(forceReturnTrue, 10*time.Millisecond)
	if res != nil {
		t.Fail()
	}
}

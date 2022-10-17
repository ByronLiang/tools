package wf

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	polling := NewPolling(15 * time.Second)
}

func fakeFileModify() {

}

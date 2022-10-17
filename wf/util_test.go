package wf

import (
	"testing"
	"time"
)

func TestCheckAndCreateFilePath(t *testing.T) {
	c := time.NewTicker(2 * time.Second)
	count := 0
loop:
	for {
		select {
		case <-c.C:
			count++
			t.Log("aa", count, time.Now().Unix())
			// the ticker handle processed time is longer than ticker duration (2s)
			// when process finished, the next process will be continue as soon
			time.Sleep(3 * time.Second)
			if count == 3 {
				break loop
			}
		}
	}
}

func TestGenWatchPath(t *testing.T) {
	count := 0
	c := time.NewTimer(2 * time.Second)
	if !c.Stop() {
		<-c.C
	}
loop:
	for {
		c.Reset(2 * time.Second)
		// when finish all process the time ticket will be reset
		// good: make sure the process time not include the duration for time ticket
		select {
		//case <-time.After(2 * time.Second):

		case <-c.C:
			count++
			t.Log("aa", count, time.Now().Unix())
			time.Sleep(3 * time.Second)
			if count == 3 {
				break loop
			}
		}
	}
}

package wf

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var testFile string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error
	// create temp file
	testFile, err = BuildTempFile("wf_*.txt")
	if err != nil {
		return
	}
	fmt.Println("build tempFile", testFile)
}

func teardown() {
	// drop temp file
	fmt.Println("drop-temp-file", testFile)
	if len(testFile) > 0 {
		os.Remove(testFile)
	}
}

func TestNew(t *testing.T) {
	var pollingCount int
	polling := NewPolling(2 * time.Second)
	defer polling.Close()
	t.Log("watch", testFile)
	err := polling.Add(testFile)
	if err != nil {
		t.Fatal(err)
	}
	go fakeFileModify(3, t)
	t1 := time.NewTimer(15 * time.Second)
loop:
	for {
		select {
		case pollingEvent := <-polling.Event:
			t.Log(pollingEvent)
			pollingCount++
		case pollingErr := <-polling.Error:
			t.Log("pollingErr: ", pollingErr)
		case <-t1.C:
			break loop
		}
	}
}

func fakeFileModify(n int, t *testing.T) {
	for i := 0; i < n; i++ {
		time.Sleep(2 * time.Second)
		if err := ioutil.WriteFile(testFile, []byte(time.Now().Format("2006-01-02 15:04:05")), 0666); err != nil {
			t.Log("write-file", err)
		}
	}
}

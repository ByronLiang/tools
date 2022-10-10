package sample

import (
	"testing"
	"time"
)

func TestNativeRefreshCounter_IncCheckReset(t *testing.T) {
	var l [5]RefreshCounter
	s1 := l[0].IncCheckReset(time.Now(), 5*time.Second)
	t.Log("s1:", s1)
	for i := 0; i < 20; i++ {
		time.Sleep(500 * time.Millisecond)
		res := l[0].IncCheckReset(time.Now(), 15*time.Second)
		t.Log(i, res)
		if res > 1 && (res-1)%12 == 0 {
			t.Log("refresh")
		}
	}
	for i := 0; i < 10; i++ {
		time.Sleep(1200 * time.Millisecond)
		res := l[0].IncCheckReset(time.Now(), 15*time.Second)
		t.Log("new", i, res)
		if res > 1 && (res-1)%12 == 0 {
			t.Log("refresh")
		}
	}
}

func TestNativeRefreshCounter_IsRefresh(t *testing.T) {
	waitTime := 80 * time.Second
	var l [5]RefreshCounter
	s1 := l[0].IsRefresh(time.Now(), waitTime, 10)
	t.Log("s1:", s1)
	for i := 0; i < 20; i++ {
		time.Sleep(300 * time.Millisecond)
		res := l[0].IsRefresh(time.Now(), waitTime, 10)
		t.Log(i, res)
		//if res > 1 && (res-1)%12 == 0 {
		//	t.Log("refresh")
		//}
	}
	for i := 0; i < 10; i++ {
		time.Sleep(1200 * time.Millisecond)
		res := l[0].IsRefresh(time.Now(), waitTime, 10)
		t.Log("new", i, res)
		//if res > 1 && (res-1)%12 == 0 {
		//	t.Log("refresh")
		//}
	}
}

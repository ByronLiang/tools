package wf

import (
	"io/ioutil"
	"log"
	"time"
)

type AutoFileModify struct {
	path          string
	counter       int
	duration      int
	lastCheckTime time.Time
}

func InitAutoFileModify(path string, duration int) (*AutoFileModify, error) {
	err := CheckAndCreateFilePath(path)
	if err != nil {
		return nil, err
	}
	autoFileModify := &AutoFileModify{
		path:          path,
		counter:       0,
		duration:      duration,
		lastCheckTime: time.Now(),
	}
	time.AfterFunc(5*time.Second, autoFileModify.CheckHandle)
	return autoFileModify, nil
}

func (a *AutoFileModify) CheckHandle() {
	//rand.Seed(time.Now().Unix())
	if err := ioutil.WriteFile(a.path, []byte(time.Now().Format("2006-01-02 15:04:05")), 0666); err != nil {
		log.Println("write-file", err)
	} else {
		a.counter++
	}
	//d := rand.Intn(a.duration)
	time.AfterFunc(time.Duration(a.duration)*time.Second, a.CheckHandle)
}

// 定时探测行为
func (a *AutoFileModify) BeatProbe() {
	log.Println("BeatProbe")
	// 容错120 second
	if a.lastCheckTime.Add(time.Duration(120+a.duration) * time.Second).Before(time.Now()) {
		log.Println("BeatProbe overdue")
	}
}

func (a *AutoFileModify) BeatCallBack() {
	a.lastCheckTime = time.Now()
}

func (a *AutoFileModify) GetWatchTarget() string {
	return a.path
}

func (a *AutoFileModify) GetHeartBeatFile() string {
	return a.path
}

func (a *AutoFileModify) GetDuration() int {
	return a.duration
}

package wf

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type RecreateFileCheck struct {
	duration      int
	filepath      string
	lastCheckTime time.Time
}

func InitRecreateFileCheck(path string, duration int) (RecreateFileCheck, error) {
	err := CheckAndCreateFilePath(path)
	if err != nil {
		return RecreateFileCheck{}, err
	}
	recreateFileCheck := RecreateFileCheck{
		filepath:      path,
		duration:      duration,
		lastCheckTime: time.Now(),
	}
	return recreateFileCheck, nil
}

func (r *RecreateFileCheck) GetWatchTarget() string {
	return filepath.Dir(r.filepath)
}

func (r *RecreateFileCheck) GetHeartBeatFile() string {
	return r.filepath
}

func (r *RecreateFileCheck) CheckHandle() {
	tick := time.NewTicker(time.Duration(r.duration) * time.Second)
	for {
		select {
		case <-tick.C:
			if err := os.Remove(r.filepath); err != nil {
				log.Println("recreateFile remove err", err.Error())
			}
			f, err := os.Create(r.filepath)
			if err != nil {
				log.Println("recreateFile create err", err.Error())
			}
			if f.Close() != nil {
				log.Println("recreateFile close err")
			}
		}
	}
}

func (r *RecreateFileCheck) BeatProbe() {
	log.Println("BeatProbe")
	// 容错120 second
	if r.lastCheckTime.Add(time.Duration(120+r.duration) * time.Second).Before(time.Now()) {
		log.Println("BeatProbe overdue")
	}
}

func (r *RecreateFileCheck) BeatCallBack() {
	r.lastCheckTime = time.Now()
}

func (r *RecreateFileCheck) GetDuration() int {
	return r.duration
}

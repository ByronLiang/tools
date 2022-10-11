package wf

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Handle func(op fsnotify.Event)

type Wf struct {
	watcher       *fsnotify.Watcher
	heartbeat     HeartBeatCheck
	handle        Handle
	modifyCounter int
}

func NewWf() *Wf {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &Wf{
		watcher: watcher,
		handle:  defaultWfHandle,
	}
}

func (w *Wf) Close() error {
	return w.watcher.Close()
}

func (w *Wf) AddFile(name string) error {
	return w.watcher.Add(name)
}

// set heartbeat to check file watch
func (w *Wf) SetHeartBeat(beat HeartBeatCheck) error {
	err := w.watcher.Add(beat.GetWatchTarget())
	if err != nil {
		return err
	}
	w.heartbeat = beat
	go w.heartbeat.CheckHandle()
	return nil
}

func (w *Wf) SetWfHandle(handle Handle) {
	w.handle = handle
}

func (w *Wf) ToWatch() {
	if w.heartbeat == nil {
		go w.watchWithoutHeartBeat()
	} else {
		go w.watch()
	}
	return
}

func (w *Wf) watchWithoutHeartBeat() {
	for {
		select {
		case event := <-w.watcher.Events:
			w.modifyCounter++
			w.handle(event)
		case err := <-w.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (w *Wf) watch() {
	checkNotifyTick := time.Tick(time.Duration(120+w.heartbeat.GetDuration()) * time.Second)
	for {
		select {
		case <-checkNotifyTick:
			w.heartbeat.BeatProbe()
		case event := <-w.watcher.Events:
			w.modifyCounter++
			if event.Name == w.heartbeat.GetHeartBeatFile() {
				w.heartbeat.BeatCallBack()
			} else {
				// none heartbeat event
				w.handle(event)
			}
		case err := <-w.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (w *Wf) GetCounter() int {
	return w.modifyCounter
}

func defaultWfHandle(event fsnotify.Event) {
	if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
		log.Println("event:", event, event.Op)
	}
}

package wf

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WfHandle func(op fsnotify.Event)

type Wf struct {
	watcher       *fsnotify.Watcher
	heartbeat     HeartBeatCheck
	handle        WfHandle
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
	err := w.watcher.Add(GetWatchTarget())
	if err != nil {
		return err
	}
	w.heartbeat = beat
	go CheckHandle()
	return nil
}

func (w *Wf) SetWfHandle(handle WfHandle) {
	w.handle = handle
}

//func (w *Wf) SetTicker(tickerFile string) error {
//    err := w.watcher.Add(tickerFile)
//    if  err != nil {
//        return err
//    }
//    w.tickerFile = tickerFile
//    w.lastGetCheckNotifyTime = time.Now()
//    go func() {
//        tickerRecreateFile(w.tickerFile, 10 * time.Second)
//        tickerRecreateFile(w.tickerFile, 30 * time.Second)
//    }()
//    return nil
//}

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
	checkNotifyTick := time.Tick(time.Duration(120+GetDuration()) * time.Second)
	for {
		select {
		case <-checkNotifyTick:
			BeatProbe()
		case event := <-w.watcher.Events:
			w.modifyCounter++
			if event.Name == GetHeartBeatFile() {
				BeatCallBack()
				log.Println("event:", event, event.Op)
			} else {
				// none heartbeat event
				w.handle(event)
			}
		case err := <-w.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

//func (w *Wf) WatchHandle() {
//    go func() {
//        checkNotifyTick := time.Tick(20 * time.Second)
//        for {
//            select {
//            case <-checkNotifyTick:
//                if w.lastGetCheckNotifyTime.Add(20 * time.Second).Before(time.Now()) {
//                    log.Println("lost event reset again")
//                }
//            case event := <-w.watcher.Events:
//                //log.Println("event:", event, event.Op)
//                if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
//                    if event.Name == w.tickerFile {
//                        //log.Println("tickerFile catch")
//                        w.lastGetCheckNotifyTime = time.Now()
//                    }
//                    //log.Println("modified file:", event.Name)
//                    w.modifyCounter++
//                }
//            case err := <-w.watcher.Errors:
//                log.Println("error:", err)
//            }
//        }
//    }()
//}

func (w *Wf) GetCounter() int {
	return w.modifyCounter
}

func defaultWfHandle(event fsnotify.Event) {
	if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
		log.Println("event:", event, event.Op)
	}
}

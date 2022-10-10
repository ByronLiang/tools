package wf

type HeartBeatCheck interface {
	GetWatchTarget() string
	GetHeartBeatFile() string
	CheckHandle()
	GetDuration() int
	BeatProbe()
	BeatCallBack()
}

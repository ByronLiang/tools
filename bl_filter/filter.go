package bl_filter

type BlFilter struct {
	// add mutex 线程安全
	Storage IStorage
}

type IStorage interface {
	Get() string
	GetAll() []string
	Save(str string, isNew bool)
}

func NewFilter(storage IStorage) *BlFilter {
	return &BlFilter{Storage: storage}
}

func (bl *BlFilter) Add(content string) {
	bl.Storage.Save(content, true)
}

func (bl *BlFilter) UpdateCurrent(content string) {
	bl.Storage.Save(content, false)
}

func (bl *BlFilter) All() []string {
	return bl.Storage.GetAll()
}

func (bl *BlFilter) Current() string {
	return bl.Storage.Get()
}

package bl_filter

import (
	"container/list"
	"sync"
)

type Queue struct {
	mu sync.RWMutex
	// 链表: 无需涉及扩容，模拟队列
	data    *list.List
	size    int
	maxSize int
}

func NewQueue(maxSize int) *Queue {
	return &Queue{
		data:    list.New(),
		maxSize: maxSize,
	}
}

func (q *Queue) GetAll() []string {
	q.mu.RLock()
	defer q.mu.RUnlock()
	res := make([]string, 0, q.data.Len())
	// 尾部链表遍历数据: 从最近数据开始遍历
	for e := q.data.Back(); e != nil; e = e.Prev() {
		if str, ok := e.Value.(string); ok {
			res = append(res, str)
		}
	}
	return res
}

func (q *Queue) Get() string {
	q.mu.RLock()
	defer q.mu.RUnlock()
	e := q.data.Back()
	if e == nil {
		return ""
	}
	if str, ok := e.Value.(string); ok {
		return str
	}
	return ""
}

func (q *Queue) Save(str string, isNew bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !isNew {
		// 链表除尾部外，其余节点已无法更新数据，因此，只能更新尾部数据
		// 取出旧数据, 并移除
		if e := q.data.Back(); e != nil {
			q.data.Remove(e)
		}
	}
	// 更新数据
	q.data.PushBack(str)
	if isNew {
		q.recollect()
	}
	return
}

func (q *Queue) recollect() {
	// 可根据队列成员数量, 将队列头部数据删除, 解决无限占用内存空间
	if q.size+1 > q.maxSize {
		if e := q.data.Front(); e != nil {
			q.data.Remove(e)
		}
		return
	}
	q.size++
}

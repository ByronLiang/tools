package singleflight

import (
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
	// 移除缓冲区数据
	forgotten bool
	dups      int
	chans     []chan<- Result
}

type Group struct {
	mu sync.Mutex            // protects m
	m  map[interface{}]*call // lazily initialized
}

type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}

func (g *Group) Do(key interface{}, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	g.doCall(c, key, fn)
	return c.val, c.err, false
}

func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)

	return ch
}

func (g *Group) doCall(c *call, key interface{}, fn func() (interface{}, error)) {
	defer func() {
		c.wg.Done()
		g.mu.Lock()
		defer g.mu.Unlock()
		if !c.forgotten {
			delete(g.m, key)
		}
		for _, ch := range c.chans {
			ch <- Result{c.val, c.err, c.dups > 0}
		}
	}()
	c.val, c.err = fn()
}

func (g *Group) Forget(key interface{}) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.forgotten = true
	}
	delete(g.m, key)
	g.mu.Unlock()
}

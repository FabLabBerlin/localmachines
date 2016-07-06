// cache running in-memory
package cache

import (
	"sync"
	"time"
)

const MAX_SECONDS = 300

var (
	mu    sync.Mutex
	items map[string]Item = make(map[string]Item)
)

func Get(key string) (o interface{}, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	i, ok := items[key]

	if !ok {
		return
	}

	if time.Since(i.Time()).Seconds() > MAX_SECONDS {
		delete(items, key)
		return nil, false
	}

	return i.o, true
}

func Invalidate(key string) {
	mu.Lock()
	defer mu.Unlock()

	delete(items, key)
}

func Put(key string, o interface{}) {
	mu.Lock()
	defer mu.Unlock()

	items[key] = Item{
		o: o,
		t: time.Now(),
	}
}

type Item struct {
	o interface{}
	t time.Time
}

func (i Item) Object() interface{} {
	return i.o
}

func (i Item) Time() time.Time {
	return i.t
}

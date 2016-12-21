package tsmap

import (
	"sync"
	"fmt"
	"time"
)

type MapElement struct {
	expiresAt uint64
	value     string
	sync.RWMutex
}

func (e *MapElement) IsExpired() bool {
	return time.Now().Unix() > int64(e.expiresAt)
}

func (e *MapElement) Update(value string, ttl uint64) error {
	e.Lock()
	defer e.Unlock()
	e.value = value
	e.expiresAt = uint64(time.Now().Unix()) + ttl
	return nil
}

type ThreadSafeMap struct {
	storage map[string]*MapElement
	sync.RWMutex
	ttl     uint64
}

func NewThreadSafeMap(ttl uint64, keys []string) *ThreadSafeMap {
	res := new(ThreadSafeMap)
	res.ttl = ttl
	res.storage = make(map[string]*MapElement, len(keys))
	for _, key := range keys {
		res.storage[key] = &MapElement{}
	}
	return res
}

type NoSuchKey struct {
	key string
}

func (e NoSuchKey) Error() string {
	return fmt.Sprintf("No such key: %s", e.key)
}

type ValueExpired struct {
	key string
}

func (e ValueExpired) Error() string {
	return fmt.Sprintf("Value for key: '%s' is expired", e.key)
}

func (m *ThreadSafeMap) Get(key string) (string, error) {

	m.RLock()
	defer m.RUnlock()

	var v *MapElement

	if v, _ = m.storage[key]; v == nil {
		return "", NoSuchKey{key}
	}

	v.RLock()
	defer v.RUnlock()
	if v.IsExpired() {
		return "", ValueExpired{key}
	}

	return v.value, nil
}

func (m *ThreadSafeMap) Set(key string, value string, ttl uint64) error {
	if ttl == 0 {
		ttl = m.ttl
	}

	m.RLock()
	defer m.RUnlock()

	var v *MapElement

	if v, _ = m.storage[key]; v == nil {
		m.RUnlock()

		m.Lock()
		if v, _ = m.storage[key]; v == nil {
			m.storage[key] = &MapElement{}
		}
		m.Unlock()

		m.RLock()
	}

	return m.storage[key].Update(value, ttl)
}
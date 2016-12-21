package tsmap

import (
	"testing"
	"runtime"
	"math/rand"
)

func TestConstructor(t *testing.T) {
	m := NewThreadSafeMap(300, []string{})
	if m == nil {
		t.Errorf("nil returned with empty key slice")
	}

	m = NewThreadSafeMap(300, []string{"some", "keys", "here"})
	if m == nil {
		t.Errorf("nil returned with non-empty key slice")
	}

}

func TestExpired(t *testing.T) {
	me := new(MapElement)
	me.Update("some val", 31)
	if me.IsExpired() {
		t.Errorf("Something wrong. %v", me)
	}
}

func TestMapElement_Update(t *testing.T) {
	me := new(MapElement)
	me.Update("some val", 0)

	if me.value != "some val" {
		t.Errorf("Something wrong. %v", me)
	}
}

func TestNewThreadSafeMap(t *testing.T) {
	m := NewThreadSafeMap(300, []string{})
	m.Set("foo", "bar", 200)
	if v, e := m.Get("foo"); v != "bar" {
		t.Errorf("Got wrong value. %v", v)
	} else if e != nil {
		t.Errorf("Got non-nil error. %v", e)
	}
}

var concurrentCacheUsageData = []struct{ k, v string }{
	{"foo", "bar"},
	{"wtf", "dddd"},
	{"k", "val"},
}

func TestConcurrentCacheUsage(t *testing.T) {

	cpus := runtime.NumCPU()
	var probablityOfWrite float32 = 0.2

	m := NewThreadSafeMap(300, []string{})

	N := 1000

	for i := 0; i < cpus * 10; i++ {
		go func() {
			for i := 0; i < N; i++ {
				if rand.Float32() < probablityOfWrite {
					k := concurrentCacheUsageData[rand.Intn(len(concurrentCacheUsageData))].k
					v := concurrentCacheUsageData[rand.Intn(len(concurrentCacheUsageData))].v
					m.Set(k, v, 0)
				} else {
					m.Get(concurrentCacheUsageData[rand.Intn(len(concurrentCacheUsageData))].k)
				}
			}
		}()
	}
}


func getMapData() ([]string, *ThreadSafeMap){
	keys := []string{"some", "key", "here"}
	m := NewThreadSafeMap(300, keys)
	return keys, m
}


func benchmarkThreadSafeMap(b *testing.B) {
	keys, m := getMapData()

	for n := 0; n < b.N; n++ {
		m.Get(keys[rand.Intn(len(keys))])
	}
}

func BenchmarkThreadSafeMap_Get(b *testing.B) {
	benchmarkThreadSafeMap(b)
}

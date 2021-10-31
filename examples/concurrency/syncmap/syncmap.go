package main

import (
	"fmt"
	"sync"
)

// go run  --race syncmap.go  can check race condition
func main() {
	// unsafeWite()
	safeWrite()
}

// will  happend error : fatal error: concurrent map writes
func unsafeWite() {
	var m = make(map[int]int)
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			m[i] = i
			fmt.Printf("m = %+v\n", m)
		}()

	}
}

type SafeMap struct {
	m map[int]int
	sync.Mutex
}

func (sm *SafeMap) Read(k int) (int, bool) {
	sm.Lock()
	defer sm.Unlock()
	result, ok := sm.m[k]
	return result, ok
}

func (sm *SafeMap) Write(k, v int) {
	sm.Lock()
	defer sm.Unlock()
	sm.m[k] = v
}

func safeWrite() {
	s := SafeMap{
		m:     map[int]int{},
		Mutex: sync.Mutex{},
	}
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			s.Write(i, i)
		}()
	}
}

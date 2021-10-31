package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	s := NewSliceNum()
	once.Do(func() {
		s.Add(16)
	})
	once.Do(func() {
		s.Add(16)
	})
	once.Do(func() {
		s.Add(16)
	})

	fmt.Printf("s = %+v\n", s)
}

type SliceNum []int

func NewSliceNum() SliceNum {
	return make(SliceNum, 0)
}

func (s *SliceNum) Add(e int) *SliceNum {
	*s = append(*s, e)
	fmt.Println("add element:", e)
	return s
}

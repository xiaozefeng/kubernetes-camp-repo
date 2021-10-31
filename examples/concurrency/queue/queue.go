package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	q := Queue{
		queue: []string{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}
	go func() {
		for {
			q.Enqueue("a")
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		q.Dequeue()
		time.Sleep(time.Second * 1)
	}

}

type Queue struct {
	queue []string
	cond  *sync.Cond
}

func (q *Queue) Enqueue(item string) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.queue = append(q.queue, item)
	fmt.Printf("putting %s to queue , notify all\n", item)
	q.cond.Broadcast()
}

func (q *Queue) Dequeue() string {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if len(q.queue) == 0 {
		fmt.Println("no data available, wait...")
		q.cond.Wait()

	}
	ret := q.queue[0]
	q.queue = q.queue[1:]
	return ret
}

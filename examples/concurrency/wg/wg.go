package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("wait by channel")
	WaitByChannel()
	fmt.Println("")
	fmt.Println("---------")

	WaitByWaitGroup()
	fmt.Println("")
	fmt.Println("---------")
}

func WaitByChannel() {
	c := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			fmt.Printf("%d ", i)
			c <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-c
	}

}

func WaitByWaitGroup() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%d ", i)
		}()
	}

	wg.Wait()
}

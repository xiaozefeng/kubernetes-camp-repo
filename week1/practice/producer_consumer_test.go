package practice_test

import (
	"log"
	"testing"
	"time"
)

func Test_1(t *testing.T) {
	var c = make(chan int)

	go func() {
		var i = 0
		for {
			i++
			log.Printf("send val:%d into channel", i)
			time.Sleep(time.Second)
			c <- i
		}
	}()

	for v := range c {
		log.Printf("received val :%d from channel\n", v)
	}
}

func Test_2(t *testing.T) {

	var c = make(chan int)
	var quit = make(chan struct{})

	go func() {
		var i = 0
		for {
			i++
			select {
			case c <- i:
				log.Printf("send val:%d into channel", i)
				time.Sleep(time.Second)
			case <-quit:
				close(c)
				log.Print("closed channel c\n")
				return
			}
		}
	}()

	for v := range c {
		log.Printf("received val :%d from channel\n", v)
		if v > 3 {
			quit <- struct{}{}
		}
	}
}

type producer struct {
	c    chan int
	quit chan struct{}
}

func (p *producer) Close() error {
	p.quit <- struct{}{}
	return nil
}

func Test_3(t *testing.T) {
	p := &producer{
		c:    make(chan int),
		quit: make(chan struct{}),
	}
	go func() {
		var i = 0
		for {
			i++
			select {
			case p.c <- i:
				log.Printf("send val:%d into channel", i)
				time.Sleep(time.Second)
			case <-p.quit:
				log.Print("closed channel c\n")
				close(p.c)
				return
			}
		}
	}()

	for v := range p.c {

		log.Printf("received val :%d from channel\n", v)
		if v > 3 {
			if err := p.Close(); err != nil {
				log.Printf("unexpectd error: %v\n", err)
			}
		}

	}
}

func Test_4(t *testing.T) {
	var producer = func(c chan int) {
		var i = 0
		for {
			i++
			log.Printf("send val:%d into channel", i)
			c <- i
			time.Sleep(time.Second)
		}
	}
	var consumer = func(c chan int) {
		for v := range c {
			log.Printf("received val :%d from channel\n", v)
		}
	}
	var c = make(chan int, 10)
	go producer(c)
	go consumer(c)

	select {}
}

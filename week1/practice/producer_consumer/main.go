package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	var c = make(chan int, 10)
	// timeout:= time.After(time.Second * 5)
	ctx, _:= context.WithTimeout(context.Background(), time.Second*10)
	for {
		select {
		case val := <-c:
			log.Printf("received val from channel : %v\n", val)
		case <-time.After(time.Second):
			var x = rand.Intn(100)
			log.Printf("send val %d into channel", x)
			c <- x
		case <-ctx.Done():
			log.Println("timeout !!! ")
			return
		}
	}
}

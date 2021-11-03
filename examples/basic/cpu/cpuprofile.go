package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("./cpuprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var result int
	for i := 0; i < 10000000; i++ {
		result += i
	}
	fmt.Printf("result = %+v\n", result)

}

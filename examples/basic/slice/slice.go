package main

import "fmt"

func main() {
	// create slice way
	createSlice()

}

func createSlice() {
	// method one
	var s1 []int // s1 == nil
	fmt.Printf("s1 = %+v\n", s1)

	// method two
	var s2 = make([]int, 0) // s2 = empty slice
	fmt.Printf("s2 = %+v\n", s2)

	// method three
	var s3 = []int{1, 2, 3}
	fmt.Printf("s3 = %+v\n", s3)
}

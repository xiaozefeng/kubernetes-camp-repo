package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string `json:"name"`
}

func main() {
	p := Person{Name: "jack"}
	// get type
	t := reflect.TypeOf(p)
	fmt.Printf("t = %+v\n", t)
	// get field
	name := t.Field(0)
	fmt.Printf("name = %+v\n", name)
	// get tag
	tag := name.Tag.Get("json")
	fmt.Println("tag:", tag)

	fmt.Println("--------------")
	t2 := T2{
		F2: "f2",
		T1: T1{
			F1: "f1",
		},
	}
	fmt.Printf("t2 = %+v\n", t2)
}

type T1 struct {
	F1 string
}

type T2 struct {
	F2 string
	T1
}

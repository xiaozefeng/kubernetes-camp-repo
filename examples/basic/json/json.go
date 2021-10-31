package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Person struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

type Address struct {
	City string `json:""`
}

func main() {
	var p = Person{
		Name: "jack",
		Address: Address{
			City: "GZ",
		},
	}
	p0, err := json.Marshal(&p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("p0 = %s\n", p0)

	var data = `{"name":"jack","address":{"City":"GZ"}}`

	var p1 Person
	err = json.Unmarshal([]byte(data), &p1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("p1 = %+v\n", p1)

	// json decode

	var p2 Person
	decoder := json.NewDecoder(strings.NewReader(data))
	err = decoder.Decode(&p2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("p2 = %+v\n", p2)

}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		value,ok := values["k"]
		fmt.Fprintln(rw, value, ok, os.Getenv(value[0]))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

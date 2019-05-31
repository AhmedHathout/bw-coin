package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	cPort = "4390"
)

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am root"))
}

func main() {
	http.HandleFunc("/", root)
	fmt.Println("start listening on port : " + cPort)
	if err := http.ListenAndServe(":"+cPort, nil); err != nil {
		log.Fatal(err)
	}
}

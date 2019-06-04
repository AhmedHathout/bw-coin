package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	cPort      = "3030"
	cBlockPath = "./res/button_block.json"
)

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bw-coin slack made by @elhmn"))
}

func bwcoinsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	text := r.FormValue("text")
	if text == "" {
		text = "No data sent !"
	}

	block, err := ioutil.ReadFile(cBlockPath)
	if err != nil {
		http.Error(w, "Error : "+err.Error(), http.StatusBadRequest)
		return
	}

	data := "{" + "\"blocks\"" + ":" + string(block) + "}"

	w.Write([]byte(data))
}

func answerInteractive(actions interface{}, url interface{}) {
	time.Sleep(5 * time.Second)
	var buf bytes.Buffer
	message := struct {
		Text string `json:"text"`
	}{Text: "Message received"}

	json.NewEncoder(&buf).Encode(message)
	http.Post(url.(string), "application/json", &buf)
}

func interactionsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	text := r.FormValue("text")
	if text == "" {
		text = "No data sent !"
	}

	payload := r.Form.Get("payload")

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//start response go routine
	{
		actions := data["actions"]
		url := data["response_url"]
		go answerInteractive(actions, url)
	}
}

func run() {
	fmt.Println("Server started at port :", cPort)
	http.HandleFunc("/", root)
	http.HandleFunc("/bwcoins", bwcoinsHandler)
	http.HandleFunc("/interactions", interactionsHandler)
	if err := http.ListenAndServe(":"+cPort, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	run()
}

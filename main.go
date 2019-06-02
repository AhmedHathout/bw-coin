package main

import (
	"fmt"
	// 	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	// 	"reflect"
)

const (
	cPort = "4390"
)

type sHeader struct {
	event     string
	delivery  string
	signature string
	cType     string
	method    string
}

type sUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

type sPullRequest struct {
	URL   string `json:"url"`
	State string `json:"state"`
	User  sUser  `json:"user"`
}

type sComment struct {
	URL   string `json:"url"`
	State string `json:"state"`
	User  sUser  `json:"user"`
}

type sReview struct {
	URL            string `json:"url"`
	State          string `json:"state"`
	User           sUser  `json:"user"`
	PullRequestURL string `json:"pull_request_url"`
}

type sRepo struct {
	HTMLURL string `json:"html_url"`
}

type sPayload struct {
	Action      string       `json:"action"`
	Review      sReview      `json:"review"`
	Comment     sComment     `json:"comment"`
	Repo        sRepo        `json:"repository"`
	PullRequest sPullRequest `json:"pull_request"`
}

func dumpRequestInfo(r *http.Request) {
	if r == nil {
		return
	}

	fmt.Println("Method : " + r.Method)
	fmt.Println("Content-type : " + r.Header.Get("content-type"))
	fmt.Println("X-GitHub-Event :" + r.Header.Get("X-GitHub-Event"))
	fmt.Println("X-GitHub-Delivery :" + r.Header.Get("X-GitHub-Delivery"))
	fmt.Println("X-Hub-Signature :" + r.Header.Get("X-Hub-Signature"))
}

func getHeaderData(r *http.Request) *sHeader {
	if r == nil {
		return nil
	}

	return &sHeader{
		method:    r.Method,
		cType:     r.Header.Get("content-type"),
		event:     r.Header.Get("X-GitHub-Event"),
		delivery:  r.Header.Get("X-GitHub-Delivery"),
		signature: r.Header.Get("X-Hub-Signature"),
	}
}

func getRequestJSON(r *http.Request) string {
	// 	defer r.Body.Close()
	return "show response body"
}

func getRequestFormData(w http.ResponseWriter, r *http.Request) sPayload {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var payloadData sPayload
	payload := r.Form.Get("payload")

	if err := json.Unmarshal([]byte(payload), &payloadData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	return payloadData
}

func dumpFormData(data sPayload) {
	fmt.Println("action: ", data.Action)
	// 	fmt.Println("review: ", reflect.ValueOf(data.Review).String())
	// 	fmt.Println("comment: ", reflect.ValueOf(data.Comment).String())
	// 	fmt.Println("repository: ", reflect.ValueOf(data.Repo).String())
}

func root(w http.ResponseWriter, r *http.Request) {

	dumpRequestInfo(r)
	headerData := getHeaderData(r)
	switch r.Header.Get("content-type") {
	case "application/json":
		str := getRequestJSON(r)
		fmt.Println(str)
	case "application/x-www-form-urlencoded":
		data := getRequestFormData(w, r)
		dumpFormData(data)
		//handle request
		{
			if headerData.event == "pull_request_review_comment" {
				switch data.Action {
				case "created":
					//check if comment creator has already left a comment on this PR
					// or pull_request state is open
					// then do nothing
					// else then
					// set a point event field to the collection (eventType, pr_url, userlogin, userid)
					// edit user point field and add corresponding point to the event
					fmt.Println("created")
				}
			}
			if headerData.event == "pull_request_review" {
				switch data.Action {
				case "submitted":
					state := data.Review.State
					fmt.Println("state : ", state)
					switch state {
					case "approved":
						fmt.Println("state == approved")
					case "commented":
						fmt.Println("state == commented")
						handleSubmittedApprovedState(state, data)
					case "request_changes":
						fmt.Println("state == request_changes")
					}
					//check review state :
					// if state === approved
					// set point to 2
					// if state === requestchanges
					// set point to 1
					// if state === pending do nothing
					// 					fmt.Println("review: ", field)
					fmt.Println("submitted")
				}
			}
		}
	}
}

func main() {
	// 	handleSubmittedApprovedState("", sPayload{})
	http.HandleFunc("/", root)
	fmt.Println("start listening on port : " + cPort)
	if err := http.ListenAndServe(":"+cPort, nil); err != nil {
		log.Fatal(err)
	}
}

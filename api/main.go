package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

const (
	cPort = "8080"
)

type sHeader struct {
	method string
}

var gSession *mgo.Session

const (
	cDBName     = "bw-coin"
	cDBUsername = "elhmn"
	cDBPassword = "mongobeti" //Later fetch password using os.Getenv()
)

type sUserDB struct {
	Login     			string `json:"login"`
	URL       			string `json:"url"`
	AvatarURL 			string `json:"avatar_url"`
	Type      			string `json:"type"`
	ID        			int    `json:"id"`
	Coins     			int    `json:"coins"`
	PullRequestsURLs	string `json:"pull_requests_urls"`
	CommentsURLs		string `json:"comments_urls"`
	// I think there is no need to keep the urls of PRs and comments but just the count
}

func createDatabase() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost"},
		Database: cDBName,
		Username: cDBUsername,
		Password: cDBPassword,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bw-coin api made by @elhmn"))
}

func getUsersData() []sUserDB {
	usersCollection := gSession.DB(cDBName).C("users")

	// Check if user already exist
	var userField []sUserDB
	if err := usersCollection.Find(nil).All(&userField); err != nil {
		panic(err)
	}
	return userField
}

// This function can be merged with the previous one. We can
// use (query interface) instead of (username string) or just use varargs
func getUserByUsername(username string) []sUserDB {
	usersCollection := gSession.DB(cDBName).C("users")

	var userField []sUserDB
	if err := usersCollection.Find(bson.M{"login": bson.RegEx{Pattern: "username"}}).All(&userField); err != nil {
		panic(err)
	}
	return userField
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		// Get username value if it exists
		searchParameter, ok := r.URL.Query()["search"]

		var (
			data []byte
			err error
		)

		if !ok || len(searchParameter) < 1 {
			data, err = json.Marshal(getUsersData())
		} else {
			data, err = json.Marshal(getUserByUsername(searchParameter[0]))
		}
		if err != nil {
			panic(err)
		}
		w.Write([]byte(data))
	}
}

func run() {
	fmt.Println("Server started at port :", cPort)
	if gSession == nil {
		fmt.Println("Setting up Database session...")
		gSession = createDatabase()
		fmt.Println("Database successfully setup !")
	}
	http.HandleFunc("/", root)
	http.HandleFunc("/users", userHandler)
	if err := http.ListenAndServe(":"+cPort, nil); err != nil {
		log.Fatal(err)
	}
	if gSession != nil {
		gSession.Close()
	}
}

func main() {
	run()
}

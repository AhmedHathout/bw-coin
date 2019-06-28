package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
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
	Login     string `json:"login"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
	ID        int    `json:"id"`
	Coins     int    `json:"coins"`
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

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		data, err := json.Marshal(getUsersData())
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

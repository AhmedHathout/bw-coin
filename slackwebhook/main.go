package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	cPort      = "3030"
	cBlockPath = "./res/listelement.json"
)

type sUserDB struct {
	Login     string
	URL       string
	AvatarURL string
	Type      string
	ID        int
	Coins     int
}

const (
	cDBName     = "bw-coin"
	cDBUsername = "elhmn"
	cDBPassword = "mongobeti" //Later fetch password using os.Getenv()
)

var gSession *mgo.Session

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

func getUserData() []sUserDB {

	usersCollection := gSession.DB(cDBName).C("users")
	var users []sUserDB
	if err := usersCollection.Find(nil).Sort("-coins").All(&users); err != nil {
		panic(err)
	}
	return users
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bw-coin slack made by @elhmn"))
}

func bwcoinHandler(w http.ResponseWriter, r *http.Request) {
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

	users := getUserData()
	blocks := ""
	for i, user := range users {
		tmp := strings.ReplaceAll(string(block), "[[login]]", user.Login)
		tmp = strings.ReplaceAll(string(tmp), "[[login_url]]", user.URL)
		tmp = strings.ReplaceAll(string(tmp), "[[avatar_url]]", user.AvatarURL)
		tmp = strings.ReplaceAll(string(tmp), "[[coins]]", strconv.Itoa(user.Coins))
		tmp = strings.ReplaceAll(string(tmp), "[[rank]]", strconv.Itoa(i+1))
		blocks += tmp
		if len(users)-1 != i {
			blocks += ","
		}
	}

	data := "{" + "\"blocks\"" + ":[" + string(blocks) + "]}"
	fmt.Println(data)

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

func interactionHandler(w http.ResponseWriter, r *http.Request) {
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
	if gSession == nil {
		fmt.Println("Setting up Database session...")
		gSession = createDatabase()
		fmt.Println("Database successfully setup !")
	}
	http.HandleFunc("/", root)
	http.HandleFunc("/bwcoin", bwcoinHandler)
	http.HandleFunc("/interaction", interactionHandler)
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

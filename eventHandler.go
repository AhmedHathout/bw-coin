package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type sUserDB struct {
	Login string
	URL   string
	Type  string
	ID    int
	Coins int
}

type sRecord struct {
	Type      string
	UserLogin string
	PrURL     string
}

const (
	cDBName          = "bw-coin"
	cDBUsername      = "elhmn"
	cDBPassword      = "mongobeti" //Later fetch password using os.Getenv()
	cInitialCoins    = 5
	cCoinsForComment = 1
)

func hasOtherCommentRecords(session *mgo.Session, data sPayload) bool {
	recordsCollection := session.DB(cDBName).C("records")

	var record []sRecord
	if err := recordsCollection.Find(bson.M{"type": "comment",
		"userlogin": data.Review.User.Login, "prurl": data.PullRequest.URL}).All(&record); err != nil {
		panic(err)
	}
	if record == nil || len(record) <= 0 {
		return false
	}
	return true
}

func dumpPayload(data sPayload) {
	fmt.Println("data:", data)
	fmt.Println("User Review url:", data.Review.URL)
	fmt.Println("User Review state:", data.Review.State)
	fmt.Println("User login:", data.Review.User.Login)
	fmt.Println("User PR URL:", data.PullRequest.URL)
	fmt.Println("User id:", data.Review.User.ID)
	fmt.Println("User url:", data.Review.User.URL)
}

func addNewCommentRecord(session *mgo.Session, data sPayload) {
	recordsCollection := session.DB(cDBName).C("records")

	if err := recordsCollection.Insert(sRecord{"comment",
		data.Review.User.Login, data.PullRequest.URL}); err != nil {
		panic(err)
	}
	fmt.Println("Comment record successfully created...")
}

func handleSubmittedApprovedState(state string, data sPayload) {
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
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	usersCollection := session.DB(cDBName).C("users")

	// Create user
	{
		// Check if user already exist
		var userField []sUserDB
		if err := usersCollection.Find(bson.M{"login": data.Review.User.Login}).All(&userField); err != nil {
			panic(err)
		}
		//Create user
		if userField == nil || len(userField) <= 0 {
			if err := usersCollection.Insert(sUserDB{data.Review.User.Login,
				data.Review.User.URL,
				data.Review.User.Type,
				data.Review.User.ID,
				cInitialCoins}); err != nil {
				panic(err)
			}
			fmt.Println("User : [", data.Review.User.Login, "] successfully created...")
		}
	}

	//Add coins for comment
	{
		if !hasOtherCommentRecords(session, data) {
			//Add coin for comment
			usersCollection.Update(bson.M{"login": data.Review.User.Login},
				bson.M{"$inc": bson.M{"coins": 1}})
			//Add new comment record
			addNewCommentRecord(session, data)
		}
	}

	// 	dumpPayload(data) // Debug
}

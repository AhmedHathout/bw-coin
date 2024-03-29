package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type sUserDB struct {
	Login     string
	URL       string
	AvatarURL string
	Type      string
	ID        int
	Coins     int
}

type sRecord struct {
	Type      string
	UserLogin string
	PrURL     string
}

var gSession *mgo.Session

const (
	cDBName            = "bw-coin"
	cDBUsername        = "elhmn"
	cDBPassword        = "mongobeti" //Later fetch password using os.Getenv()
	cInitialCoins      = 5
	cCoinsForComment   = 1
	cCoinsForApproval  = 2
	cCoinsForChanges   = 1
	cMaxCoinsPerPR     = 2
	cMaxCoinsLostPerPR = 4
)

const (
	eRecordCommented = "commented"
	eRecordApproved  = "approved"
	eRecordChanges   = "changes_requested"
)

func typeToCoins(recordType string) int {
	coins := 0

	switch recordType {
	case eRecordApproved:
		coins = cCoinsForApproval
	case eRecordCommented:
		coins = cCoinsForComment
	case eRecordChanges:
		coins = cCoinsForChanges
	}
	return coins
}

//Total coins a user collected on a PR
func totalCoinsPerPR(session *mgo.Session, data sPayload) int {
	recordsCollection := session.DB(cDBName).C("records")
	coins := 0

	var records []sRecord
	if err := recordsCollection.Find(bson.M{"userlogin": data.Review.User.Login,
		"prurl": data.PullRequest.URL}).All(&records); err != nil {
		panic(err)
	}
	for _, e := range records {
		coins += typeToCoins(e.Type)
	}
	return coins
}

func totalCoinsLostPerPR(session *mgo.Session, data sPayload) int {
	recordsCollection := session.DB(cDBName).C("records")
	coins := 0

	var records []sRecord
	if err := recordsCollection.Find(bson.M{"prurl": data.PullRequest.URL}).All(&records); err != nil {
		panic(err)
	}
	for _, e := range records {
		coins += typeToCoins(e.Type)
	}
	return coins
}

func hasOtherRecords(session *mgo.Session, data sPayload, recordType string) bool {
	recordsCollection := session.DB(cDBName).C("records")

	var records []sRecord
	if err := recordsCollection.Find(bson.M{"type": recordType,
		"userlogin": data.Review.User.Login,
		"prurl":     data.PullRequest.URL}).All(&records); err != nil {
		panic(err)
	}
	if records == nil || len(records) <= 0 {
		return false
	}
	return true
}

func dumpPayload(data sPayload) {
	fmt.Println("data:", data)
	fmt.Println("User Review url:", data.Review.URL)
	fmt.Println("User Review state:", data.Review.State)
	fmt.Println("User login:", data.Review.User.Login)
	fmt.Println("User avatar_url:", data.Review.User.AvatarURL)
	fmt.Println("User PR URL:", data.PullRequest.URL)
	fmt.Println("User id:", data.Review.User.ID)
	fmt.Println("User url:", data.Review.User.URL)
	fmt.Println("User PR id:", data.PullRequest.User.ID)
	fmt.Println("User PR url:", data.PullRequest.User.Login)
	fmt.Println("User PR avatar_url:", data.PullRequest.User.AvatarURL)
}

func addNewRecord(session *mgo.Session, data sPayload, recordType string) {
	recordsCollection := session.DB(cDBName).C("records")

	if err := recordsCollection.Insert(sRecord{recordType,
		data.Review.User.Login, data.PullRequest.URL}); err != nil {
		panic(err)
	}
	fmt.Println("Record :", recordType, " : successfully created !")
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

// Create user
func createUserIfDoesNotExist(usersCollection *mgo.Collection, data sPayload) {
	// Check if user already exist
	var userField []sUserDB
	if err := usersCollection.Find(bson.M{"login": data.Review.User.Login}).All(&userField); err != nil {
		panic(err)
	}
	//Create user
	if userField == nil || len(userField) <= 0 {
		if err := usersCollection.Insert(sUserDB{data.Review.User.Login,
			data.Review.User.URL,
			data.Review.User.AvatarURL,
			data.Review.User.Type,
			data.Review.User.ID,
			cInitialCoins}); err != nil {
			panic(err)
		}
		fmt.Println("User : [", data.Review.User.Login, "] successfully created !")
	}
}

func getIncrement(coins int, gain int) int {
	inc := gain
	diff := cMaxCoinsPerPR - coins - inc
	if diff < 0 {
		if inc = cMaxCoinsPerPR - coins; inc < 0 {
			inc = 0
		}
	}
	return inc
}

func getDecrement(coins int, gain int) int {
	inc := gain
	diff := cMaxCoinsLostPerPR - coins - inc
	if diff < 0 {
		if inc = cMaxCoinsLostPerPR - coins; inc < 0 {
			inc = 0
		}
	}
	return -inc
}

func handleSubmittedState(data sPayload) {
	state := data.Review.State
	usersCollection := gSession.DB(cDBName).C("users")
	createUserIfDoesNotExist(usersCollection, data)

	if data.PullRequest.User.Login == data.Review.User.Login {
		fmt.Println("Can't handle coins as ", data.PullRequest.User.Login, " owns this PR...")
		return
	}

	if hasOtherRecords(gSession, data, state) {
		fmt.Println("Can't handle", state, " state : Request already exist !")
		return
	}

	//Add coins
	coins := totalCoinsPerPR(gSession, data)
	inc := getIncrement(coins, typeToCoins(state))
	usersCollection.Update(bson.M{"login": data.Review.User.Login},
		bson.M{"$inc": bson.M{"coins": inc}})

	//Remove coins
	coins = totalCoinsLostPerPR(gSession, data)
	inc = getDecrement(coins, typeToCoins(state))
	usersCollection.Update(bson.M{"login": data.PullRequest.User.Login},
		bson.M{"$inc": bson.M{"coins": inc}})

	addNewRecord(gSession, data, state)
}

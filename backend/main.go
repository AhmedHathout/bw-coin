package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type User struct {
	Login     string
	URL       string
	AvatarURL string
	Type      string
	ID        int
	Coins     int
}

func connectToMongo() *mongo.Client{

	// TODO Get the right mongodb URI
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func getAllUsers(client *mongo.Client) []*User{
	userCollection := client.Database("bw-coin").Collection("users")

	// Find all users
	usersCursor, err := userCollection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		panic(err)
	}

	// Convert from Cursor to []*User
	var users []*User
	for usersCursor.Next(context.TODO()){
		var user User

		err :=  usersCursor.Decode(&user)

		if err != nil {
			panic(err)
		}

		users = append(users, &user)
	}

	if err := usersCursor.Err(); err != nil {
		panic(err)
	}

	usersCursor.Close(context.TODO())

	return users
}

func closeConnectionToMongoDB(client *mongo.Client) {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func main() {

	// TODO May need to change the port. Docker?
	http.ListenAndServe(":3001", nil)
}

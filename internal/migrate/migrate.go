package migrate

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
)

func Up(dbName string, db *mongo.Client) {
	f, err := ReadFromFile("users_go.json")
	if err != nil {
		log.Fatal(err)
	}
	var users map[string][]interface{}
	err = json.Unmarshal([]byte(f), &users)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.Database(dbName).Collection("users")
	_, err = collection.InsertMany(context.TODO(), users["objects"])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data uploaded to database!")
}

func ReadFromFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

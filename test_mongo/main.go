package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
)

type AlertMessage struct {
	date    time.Time
	id      int64
	message string
}

func parseTelegramMessage(id int64, command string, message string) AlertMessage {
	tempString := strings.TrimSpace((strings.TrimLeft(message, command)))

	tempArr := strings.SplitN(tempString, " ", 2)
	if len(tempArr) != 2 {
		fmt.Println("Error")
	}
	t, _ := time.Parse(time.RFC3339, tempArr[0])
	return AlertMessage{
		date:    t,
		id:      id,
		message: tempArr[1],
	}
}

func main() {

	//client, err := mongo.NewClient("mongodb://localhost:27017 -u 'adminT' -p 'LulaLa@5678.comT' --authenticationDatabase adminT")
	client, err := mongo.NewClientWithOptions("mongodb://localhost:27017",
		&options.ClientOptions{
			ConnString: connstring.ConnString{Username: "admin", Password: "LulaLa@5678.com", Database: "admin"}})
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK!")
	}
	fmt.Println(client)

	collection := client.Database("data_alarm").Collection("test")
	fmt.Println(collection)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	testAlert := parseTelegramMessage(777, "set", "2018-12-19T18:47:00+03:00 message_t")

	res, err := collection.InsertOne(ctx, bson.M{"id": testAlert.id, "date": testAlert.date, "message": testAlert.message})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.InsertedID)
	}
}

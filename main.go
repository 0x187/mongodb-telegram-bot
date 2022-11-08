package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	Bandwidth string
	Expire    string
	Name      string
	Usage     string
	_id       string
	Data      string
}

func rplay(msg string) string {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	title := msg
	uri := ("mongodb url")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("vpndb").Collection("vpn")

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"Name", title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
	}
	if err != nil {
		return ("No data was found with the ID")
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}

	u := User{}
	err = json.Unmarshal(jsonData, &u)
	if err != nil {
		log.Fatal(err)
	}
	v := fmt.Sprintf("👤  %s  \n📶 %s\n📈  %s\n⌛  %s  \n \n ❌ last refresh date: %s ", u.Name, u.Bandwidth, u.Usage, u.Expire, u.Data)
	return v

}

func main() {
	bot, err := tgbotapi.NewBotAPI("telegram bot api token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, rplay(update.Message.Text))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

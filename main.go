package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/vallard/spark"
)

type BotConfig struct {
	Id       string // id of the bot for logging
	UserId   string // id of the user for logging.
	Token    string
	Email    string
	SparkId  string
	Commands []string
	Actions  map[string]func(string) error
	RoomId   string
}

func (b *BotConfig) Respond(input string) error {
	for _, cmd := range b.Commands {
		if strings.Contains(strings.ToLower(input),
			strings.ToLower(cmd)) {
			// we have a match!
			return b.Actions[cmd](input)
		}
	}
	return errors.New("No response for this command: " + input)
}

var build = "1"
var s *spark.Spark
var bot BotConfig

func getMessageInfo(data map[string]interface{}) spark.Message {
	var m spark.Message
	for k, v := range data {
		// make sure the value is of type string.
		if reflect.TypeOf(v) == reflect.TypeOf("") {
			vv := v.(string)
			switch k {
			case "id":
				m.Id = vv
			case "roomId":
				m.RoomId = vv
			case "roomType":
				m.RoomType = vv
			case "text":
				m.Text = vv
			case "personId":
				m.PersonId = vv
			case "personEmail":
				m.PersonEmail = vv
			case "markdown":
				m.Markdown = vv
			case "html":
				m.Html = vv
			case "created":
				tt, err := time.Parse("2006-01-02T03:04:05.000Z", vv)
				if err == nil {
					m.Created = tt
				}
			default:
				log.Printf("unknown key: %s\n", k)
			}
		}
	}
	return m
}

// note: logging and billing should only be done here if we respond.
func sendResponse(message spark.Message) error {
	// create a new message
	m, err := s.CreateMessage(message)
	if err != nil {
		log.Printf("Unable to create message.\nM: %v\n", m)
	}
	return err
}

func handleWebhook(w spark.Webhook) {
	// see if there is a message with this spark webhook.
	message := getMessageInfo(w.Data)

	// assuming we have a message from the data, see if we can get it.
	if message.Id == "" || message.RoomId == "" {
		log.Println("message had no ID or RoomID associated with it.")
		return
	}

	m, err := s.GetMessage(message.Id)
	if err != nil {
		log.Println(err)
		return
	}

	if m.PersonId == bot.SparkId {
		log.Println("Ignoring message from myself")
		return
	}

	log.Printf("Got a message from %s.  It says: %s\n", m.PersonEmail, m.Text)
	bot.RoomId = message.RoomId
	err = bot.Respond(m.Text)
}

func main() {

	// token will be validated before getting to this point as the bot needs to be registered.
	/* variables that will be subbed per each bot.  */
	bot.Token = "M2I1MGZmNWUtYWJjMy00ZjJjLTgwZmYtZjBmN2M5MGQ3MmEyYTQ1YzU3N2UtYjM0"
	bot.Commands = []string{"/help", "/code"}
	bot.Id = "5882607b0000000000000000"
	bot.UserId = "000000169293000000169293"
	bot.Email = "berlin@sparkbot.io"
	bot.SparkId = "Y2lzY29zcGFyazovL3VzL1BFT1BMRS8yYmQzNzNiZS00ODY2LTQxYzUtYTZlNC1jODBlZTU5MmM2ZjI"

	bot.Actions = map[string]func(string) error{

		"/help": f0,

		"/code": f1,
	}

	// set up our spark client.  Only want one of these.
	s = spark.New(bot.Token)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a request:\n  %v\n\n", r)
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			for {
				var wh spark.Webhook
				if err := decoder.Decode(&wh); err != nil {
					log.Print(err)
					break
				}
				// do something with the message.
				handleWebhook(wh)
			}
		}
		if r.Method == "GET" {
			fmt.Fprintf(w, "I'm alive and waiting for spark webhooks!")

		}
	})

	http.ListenAndServe(":8080", nil)
}

// action: consists of Command, Code, Text, Language

func f0(str string) error {

	newMessage := spark.Message{
		RoomId: bot.RoomId,
		Text:   "I can tell you about /berlin and /tiergarten",
	}
	return sendResponse(newMessage)

}

func f1(str string) error {

	newMessage := spark.Message{
		RoomId:   bot.RoomId,
		Markdown: "## Code\n```\ngo run main.go\n```",
	}
	return sendResponse(newMessage)

}

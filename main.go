package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
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

var pics = [...]string{
	"http://ilosm.cdnize.com/wp-content/uploads/620-rick-james-music-facts.imgcache.rev1406146254157.web_.jpg",
	"https://s-media-cache-ak0.pinimg.com/originals/b4/7c/b7/b47cb74aefd12240acc9b12519f7868a.jpg",
	"http://ilosm.cdnize.com/wp-content/uploads/rick4.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/4/4b/Rick_James_in_Lifestyles_of_the_Rich_1984.JPG",
}

var quotes = [...]string{
	"How you doing sugar?",
	"She's a very kinky girl, The kind you don't take home to mother",
	"I betcha I'll make you holler.",
	"Now is my time.  Everything I've done up to this point is just a warm up. This is where it all begins.",
	"If anything I consider myself non-violent.  I'm from the hippy era, peace, love, groovy.",
	"Get up on this funk!",
	"We're gonna dance on the funk and make love on this song.",
	"Punk Funk means to be one with yourself. To be rebellious, aggressive, able to do and say what you feel at all times, without inflicting mental or spiritual pain.",
	"I've had it all. I've done it all. I've seen it all.",
	"Funkers are people who dig the funk; Little funkers, Big funkers, Old funkers, Young funkers, Foxy funkers, Mother funkers, Papa funkers.",
}

var src = rand.NewSource(time.Now().Unix())
var r = rand.New(src)

func (b *BotConfig) Respond(input string) error {
	asdf
	newMessage := spark.Message{
		RoomId: bot.RoomId,
		//Files:  []string{"http://ilosm.cdnize.com/wp-content/uploads/620-rick-james-music-facts.imgcache.rev1406146254157.web_.jpg"},
		//Text:  "Hey Baby",
		Files: []string{pics[r.Intn(len(pics))]},
		Text:  quotes[r.Intn(len(quotes))],
	}
	return sendResponse(newMessage)
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

	if m.PersonEmail == bot.Email {
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
	bot.Email = "rickjames@sparkbot.io"
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

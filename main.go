package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/vallard/spark"
)

type BotConfig struct {
	Token    string
	Keyword  string
	Response string
	RoomId   string
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

func sendHello() error {
	log.Println("Calling Send Hello")
	// create a new message
	newMessage := spark.Message{
		RoomId: bot.RoomId,
		Text:   bot.Response,
	}
	m, err := s.CreateMessage(newMessage)
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

	//log.Printf("Room Id: %s\n", bot.RoomId)
	if strings.Contains(strings.ToLower(m.Text),
		strings.ToLower(bot.Keyword)) {
		log.Println("Someone said the key word to me")
		// set the bot room ID of where we'll send our message.
		bot.RoomId = message.RoomId
		err = sendHello()
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Printf("Someone mentioned me, but didn't say the magic word: %s\n", bot.Keyword)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "sparkbot maker"
	app.Usage = "spark bot maker"
	app.Action = run
	app.Version = fmt.Sprintf("0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			Usage:  "token of spark bot",
			EnvVar: "SPARK_TOKEN",
		},
		cli.StringFlag{
			Name:   "keyword",
			Usage:  "keyword to look for",
			EnvVar: "KEYWORD",
		},
		cli.StringFlag{
			Name:   "response",
			Usage:  "response you want for when keyword is found",
			EnvVar: "RESPONSE",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	bot.Token = c.String("token")
	bot.Keyword = "hello"
	bot.Response = "hello back to you my friend!"
	if c.String("keyword") != "" {
		bot.Keyword = c.String("keyword")
	}
	if c.String("response") != "" {
		bot.Response = c.String("response")
	}

	// set up our spark client.  Only want one of these.
	s = spark.New(bot.Token)

	http.HandleFunc("/spark-hook", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a request:\n  %v\n\n", r)
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			for {
				var wh spark.Webhook
				if err := decoder.Decode(&wh); err != nil {
					break
				}
				// do something with the message.
				handleWebhook(wh)
			}
		}
		fmt.Fprintf(w, "Thanks for playing, %q\n", r.RequestURI)
	})

	return http.ListenAndServe(":8080", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
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
				tt, err := time.Parse("2006-01-02T03:04:05+00:00Z", vv)
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
	// create a new spark client
	s := spark.New(bot.Token)

	// create a new message
	newMessage := spark.Message{
		RoomId: bot.RoomId,
		Text:   bot.Response,
	}
	_, err := s.CreateMessage(newMessage)
	return err
}

func handleWebhook(w spark.Webhook) {
	fmt.Printf("this is the data: %v\n", w.Data)
	// see if there is a message with this spark webhook.
	message := getMessageInfo(w.Data)

	// once we have the message information, we need to request the message contents.

	log.Printf("Handling webhook for spark bot.  Message:  %v\n", message)
	log.Printf("Message Text is:  %s\n", message.Text)
	bot.RoomId = message.RoomId
	if strings.Contains(message.Text, bot.Keyword) {
		log.Println("Someone said hello to me")
	} else {
		log.Printf("Someone mentioned me, but didn't say the magic word: %s\n", bot.Keyword)
	}
	// say hello anyway
	err := sendHello()
	if err != nil {
		log.Println(err)
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

	http.HandleFunc("/spark-hook", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a request:\n  %v\n\n", r)
		if r.Method == "POST" {
			dec := json.NewDecoder(r.Body)
			for {
				var wh spark.Webhook
				if err := dec.Decode(&wh); err == io.EOF {
					break
				} else if err != nil {
					log.Println(err)
				}
				// do something with the message.
				handleWebhook(wh)
			}
		}
		fmt.Fprintf(w, "Thanks for playing, %q\n", r.RequestURI)
	})

	return http.ListenAndServe(":8080", nil)
}

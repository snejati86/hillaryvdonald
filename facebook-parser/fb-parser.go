package main

import (
	"fmt"
	"log"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"hash/fnv"
	"github.com/streadway/amqp"
	"os"
	"time"
	"encoding/json"
)

type FacebookFeed struct {
	Url string
	Body string
	Id uint32
	Time int64
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

var channel *amqp.Channel

func parseFeed(s string){

	resp, err := http.Get(s)
	if err != nil {
		failOnError(err,"Could not get the facebook page.")
		// handle error
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()
			//fmt.Println(t)
			if strings.Contains(t.String(),"_5pbx userContent"){
				for {
					tt:=z.Next()
					switch  {
					case tt == html.StartTagToken:
						t := z.Token()
						if t.Data == "p" {
							_=z.Next()
							body:=z.Token().Data
							id:=hash(body)
							feedModel:=FacebookFeed{s,body,id,time.Now().Unix()}

							bytes,err:=json.Marshal(feedModel)
							fmt.Println(string(bytes))
							err = channel.Publish(
								"facebook",     // exchange
								"", // routing key
								false,  // mandatory
								false,  // immediate
								amqp.Publishing {
									ContentType: "text/plain",
									Body:        bytes,
								})
							failOnError(err, "Failed to publish a message")

							break;

						}
					}
					break;

				}
			}

		}
	}

}
func main() {

	rabbit:=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	url:=os.Getenv("FACEBOOK_FEED")
	if rabbit == "" || url == ""{
		panic("No environment for rabbit")
	}
	conn, err := amqp.Dial("amqp://guest:guest@"+rabbit+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	failOnError(err, "Failed to declare a queue ")



	failOnError(err, "Failed to bind to the queue.")

	err = channel.ExchangeDeclare(
		"facebook",
		"fanout",
		true,
		false,
		false,
		false,
		nil)

	failOnError(err, "Failed to create the exchange for facebook.")
	tickChan := time.NewTicker(time.Minute * 10).C
	parseFeed(url)
	for {
		select {
		case <- tickChan:
			parseFeed(url)
		}
	}

}
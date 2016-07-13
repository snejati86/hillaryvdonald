package main

import (
	"github.com/SlyMarbo/rss"
	"fmt"
	"time"
	"strings"
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


var channel *amqp.Channel
var queue amqp.Queue

func parseFeed(url string){
	    feed, err := rss.Fetch(url)
	    if err != nil {
			panic(err)
	    }else {
			for _, item := range feed.Items{
					if itemMeetCriteria(item){
						fmt.Println(item.Title)
						if itemIsNew(item) {
							b, err := json.Marshal(item)
							err = channel.Publish(
								"",     // exchange
								queue.Name, // routing key
								false,  // mandatory
								false,  // immediate
								amqp.Publishing {
									ContentType: "text/plain",
									Body:        b,
								})
							failOnError(err, "Failed to publish a message")
						}
					}

			}
		}
}

func itemMeetCriteria(item *rss.Item) bool {
	return strings.Contains(item.Title,"Trump") || strings.Contains(item.Title,"Clinton")
}

func itemIsNew(item *rss.Item) bool{

	return true
}



func main() {

	queueIp:=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	rssUrl:=os.Getenv("RSS_URL")
	if queueIp == "" || rssUrl == ""{
		failOnError(nil,"Environment not ready")
	}
	conn, err := amqp.Dial("amqp://"+"guest:guest@"+queueIp+":5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	channel=ch
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"feeds", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	queue=q
	tickChan := time.NewTicker(time.Minute * 10).C
	parseFeed(rssUrl)
	for {
		select {
		case <- tickChan:
			parseFeed(rssUrl)
		}
	}
}

package main

import (
	"github.com/SlyMarbo/rss"
	"fmt"
	"strings"
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	_ "github.com/hashicorp/golang-lru"
	"os"
	"time"
	"github.com/hashicorp/golang-lru"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

const (
	INTERVAL = 60*10
)


/**
Will I ever refactor this crap ? Probably not.
 */
var channel *amqp.Channel
var cache map[string]*rss.Item;

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
								"rss",     // exchange
								"", // routing key
								false,  // mandatory
								false,  // immediate
								amqp.Publishing {
									ContentType: "text/plain",
									Body:   b     ,
								})
							failOnError(err, "Failed to publish a message")
						}
					}

			}
		}
}

func itemMeetCriteria(item *rss.Item) bool {
	return strings.Contains(strings.ToLower(item.Title),"trump")|| strings.Contains(strings.ToLower(item.Title),"clinton")
}

var l *lru.Cache

func itemIsNew(item *rss.Item) bool{
	if cache[item.Title] == nil {

		fmt.Println("not cahced")
		l.Add(item.Title,item)
		return true;
	}else{
		fmt.Println("cached")
		return false;
	}
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
	failOnError(err, "Failed to open a channel")
	defer ch.Close()


	err = ch.ExchangeDeclare(
		"rss",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err,"Failed to open exchange.");

	channel=ch


	tickChan := time.NewTicker(time.Second * INTERVAL).C
	cache = make(map[string]*rss.Item)
	onEvicted := func(k interface{}, v interface{}) {
		str, _ := k.(string)
		delete(cache,str)
	}
	l,_=lru.NewWithEvict(100,onEvicted)
	parseFeed(rssUrl)
	for {
		select {
		case <- tickChan:
			fmt.Println("Parsing")
			parseFeed(rssUrl)
		}
	}
}

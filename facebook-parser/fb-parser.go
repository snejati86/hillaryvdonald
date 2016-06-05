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
	"github.com/gocql/gocql"
	"strconv"
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
var cache map[uint32]bool
var session *gocql.Session

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
							if _, ok := cache[id]; !ok {
								cache[id]=true
								//NOT IN THE MAP.
								err =session.Query("INSERT INTO hvd.fb (fbid,insertion_time,body,owner) VALUES (? , ? , ?,? )",strconv.Itoa(int(feedModel.Id)),gocql.TimeUUID().Timestamp(),feedModel.Body,feedModel.Url).Exec()
								failOnError(err,"Could not insert")
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
							}

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
	cache = make(map[uint32]bool)

	rabbit:=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	cassandra := os.Getenv("CASSANDRA_SERVICE_HOST")
	url:=os.Getenv("FACEBOOK_FEED")
	if rabbit == "" || url == "" || cassandra  == ""{
		panic("No environment for rabbit")
	}
	conn, err := amqp.Dial("amqp://guest:guest@"+rabbit+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	failOnError(err, "Failed to declare a queue ")


	cluster := gocql.NewCluster(cassandra)
	//cluster.Keyspace = keyspace
	session, err = cluster.CreateSession()

	if err != nil {
		failOnError(err, "Can't connect to cassandra hvd")
	}
	err =session.Query("CREATE KEYSPACE IF NOT EXISTS hvd WITH REPLICATION  = { 'class':'SimpleStrategy','replication_factor':'1'}").Exec()
	if err != nil {
		failOnError(err,"Can't create or open keyspace")
	}

	err =session.Query("CREATE TABLE IF NOT EXISTS hvd.fb  (fbid text PRIMARY KEY,  insertion_time timestamp, body text , owner text)").Exec()
	if err != nil {
		failOnError(err,"Can't create table")
	}

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
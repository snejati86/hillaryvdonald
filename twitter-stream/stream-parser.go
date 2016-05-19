package main

import (
	"fmt"
	"log"
	"os"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/streadway/amqp"
	"encoding/json"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
const(
	TWEETS_EXCHANGE_NAME = "tweets"
	TWEETS_EXCHANGE_TYPE = "fanout"
	TWEET_FILTER = []string{"@realDonaldTrump","@HillaryClinton"}
	queueUserName=os.Getenv("Q_USER")
	queuePassword=os.Getenv("Q_PASS")
	queueAddress=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessToken = os.Getenv("ACCESS_TOKEN")
	accessSecret = os.Getenv("ACCESS_SECRET")
)


func main() {

	//Get properties.


	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" || queueAddress == "" || queuePassword == "" || queueUserName == "" {
		log.Fatal("Consumer key/secret and Access token/secret required/queue Information")
	}

	conn, err := amqp.Dial("amqp://"+queueUserName+":"+queuePassword+"@"+queueAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	//TODO : Will these work in case of a Ctrl+C ?
	defer conn.Close()

	//Open a channel.
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//Declase exchange.
	err = ch.ExchangeDeclare(TWEETS_EXCHANGE_NAME,TWEETS_EXCHANGE_TYPE,true,false,false,false,nil)

	failOnError(err, "Failed to declare a queue")


	//Setup twitter connection.
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	/**
	A Changel to send and recieve tweets.
	 */
	tweetChannel := make(chan *twitter.Tweet)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tweetChannel<-tweet
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         TWEET_FILTER,
		StallWarnings: twitter.Bool(true),
	}


	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits

	/**
	It does scare me how fast the last few years passed me by. hollow,uneventful and dull like ephemeral ports.
	In a staring contest with sunny san diego.
	 */
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	for {
		select {
		case tweet:=<- tweetChannel:
			bytes,err:=json.Marshal(tweet)
			if err != nil {
				fmt.Println("Unable to marshal the tweet")
			}else{

				err:=ch.Publish(
					TWEETS_EXCHANGE_NAME,     // exchange
					"", // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing {
						ContentType: "text/plain",
						Body: bytes,
					})
				if err != nil {
					log.Fatal("%s: %s", "Unable to publish to queue", err)
				}else{
					log.Println(tweet.Text)
				}
			}
		}
	}
}
package main
import (
	"bufio"
	"os"
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"github.com/dghubble/go-twitter/twitter"
	"strings"
	"encoding/json"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type ProfanityEnvelope struct {
	Words []string
	Owner string
	Time string
	Toward string
}

func parseMap(fileName string) map[string]bool {
	set := make(map[string]bool)
	inFile, _ := os.Open(fileName)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		set[scanner.Text()]=true
	}
	return set
}

func checkProfanity( sentence string, profanityMap map[string]bool) []string{
	var found []string = nil
	for _,word := range strings.Fields(sentence) {
		if profanityMap[word] {
			found = append(found,word)
		}
	}
	return found
}

func main() {
	//set := make(map[string]bool)
	rabbit:=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	if rabbit == ""{
		panic("No environment for rabbit")
	}
	badWords :=parseMap("bad-words.txt")
	conn, err := amqp.Dial("amqp://guest:guest@"+rabbit+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"", // name
		false,   // durable
		false,   // delete when unused
		true,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue ")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"tweets", // exchange
		false,
		nil)

	failOnError(err, "Failed to bind to the queue.")

	err = ch.ExchangeDeclare(
		"profanity",
		"fanout",
		true,
		false,
		false,
		false,
		nil)

	failOnError(err, "Failed to create the exchange for profanity.")

	failOnError(err, "Failed to declare a queue ( profanity )")

	msgs, err := ch.Consume(
		q.Name, // queue
		"profanity-parser",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var tweet twitter.Tweet
			err:=json.Unmarshal(d.Body,&tweet)
			if err != nil {
				fmt.Println("Failed to parse tweet")
			}
			content :=tweet.Text
			log.Println(content)
			profanity:=checkProfanity(content,badWords)
			if profanity != nil {
				envelope:=ProfanityEnvelope{}
				if strings.Contains(content,"@realDonaldTrump"){
					envelope.Toward="Donald"
				}else{
					envelope.Toward="Hillary"
				}
				envelope.Owner=tweet.User.ScreenName
				envelope.Words=profanity
				envelope.Time=tweet.CreatedAt
				env,_:=json.Marshal(envelope)
				fmt.Println(envelope)
				err = ch.Publish(
										"profanity",     // exchange
										"", // routing key
										false,  // mandatory
										false,  // immediate
										amqp.Publishing {
											ContentType: "text/plain",
											Body:        env,
										})
									failOnError(err, "Failed to publish a message")
			}

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}


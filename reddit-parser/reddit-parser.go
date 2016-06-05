package main

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/hashicorp/golang-lru"
	"strings"
	"os"
	"github.com/streadway/amqp"
	"log"
)

const (
	SECONDS = 3
)
type Comment struct {
	Kind string `json:"kind"`
	Data struct {
		     Modhash string `json:"modhash"`
		     Children []struct {
			     Kind string `json:"kind"`
			     Data struct {
					  SubredditID string `json:"subreddit_id"`
					  LinkTitle string `json:"link_title"`
					  BannedBy interface{} `json:"banned_by"`
					  RemovalReason interface{} `json:"removal_reason"`
					  LinkID string `json:"link_id"`
					  LinkAuthor string `json:"link_author"`
					  Likes interface{} `json:"likes"`
					  Replies string `json:"replies"`
					  UserReports []interface{} `json:"user_reports"`
					  Saved bool `json:"saved"`
					  ID string `json:"id"`
					  Gilded int `json:"gilded"`
					  Archived bool `json:"archived"`
					  Stickied bool `json:"stickied"`
					  Author string `json:"author"`
					  ParentID string `json:"parent_id"`
					  Score int `json:"score"`
					  ApprovedBy interface{} `json:"approved_by"`
					  Over18 bool `json:"over_18"`
					  ReportReasons interface{} `json:"report_reasons"`
					  Controversiality int `json:"controversiality"`
					  Body string `json:"body"`
					  Edited bool `json:"edited"`
					  AuthorFlairCSSClass string `json:"author_flair_css_class"`
					  Downs int `json:"downs"`
					  BodyHTML string `json:"body_html"`
					  Quarantine bool `json:"quarantine"`
					  Subreddit string `json:"subreddit"`
					  ScoreHidden bool `json:"score_hidden"`
					  Name string `json:"name"`
					  Created float64 `json:"created"`
					  AuthorFlairText string `json:"author_flair_text"`
					  LinkURL string `json:"link_url"`
					  CreatedUtc float64 `json:"created_utc"`
					  Ups int `json:"ups"`
					  ModReports []interface{} `json:"mod_reports"`
					  NumReports interface{} `json:"num_reports"`
					  Distinguished interface{} `json:"distinguished"`
				  } `json:"data"`
		     } `json:"children"`
		     After string `json:"after"`
		     Before interface{} `json:"before"`
	     } `json:"data"`
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


func main() {
	rabbit:=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	subreddit := os.Getenv("SUBREDDIT")
	if rabbit == "" || subreddit == ""{
		panic("No environment for rabbit")
	}
	conn, err := amqp.Dial("amqp://guest:guest@"+rabbit+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"reddit",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err,"Failed to open exchange.");
	cache := make(map[string][]byte)
	onEvicted := func(k interface{}, v interface{}) {
		str, _ := k.(string)

		delete(cache,str)
	}
	l,_:=lru.NewWithEvict(100,onEvicted)

	uploadTime := time.NewTicker(time.Second * SECONDS).C
	for {
		select{
		case <-uploadTime:
			go func (){
				resp, err := http.Get("https://www.reddit.com/r/"+subreddit+"/comments/.json?limit=100")
				if err != nil {
					fmt.Println("Error")
				}else{
					defer resp.Body.Close()
					comment:=Comment{}
					body, _ := ioutil.ReadAll(resp.Body)
					err:=json.Unmarshal(body,&comment)
					if err != nil {
						fmt.Println(err)
					}else{


						for _,v:=range comment.Data.Children{
							//fmt.Println(v.Data.Body)
							if strings.Contains(strings.ToLower(v.Data.Body),"trump") || strings.Contains(strings.ToLower(v.Data.Body),"clinton") || strings.Contains(strings.ToLower(v.Data.Body),"hillary") {
								if cache[v.Data.ID] == nil {
									fmt.Println("Not cached")
									cache[v.Data.ID] = body
									fmt.Println(v.Data.Body)
									l.Add(v.Data.ID,body)
									msg,err:=json.Marshal(v)
									if err != nil{
										continue;
									}
									err = ch.Publish(
										"reddit",     // exchange
										"", // routing key
										false,  // mandatory
										false,  // immediate
										amqp.Publishing {
											ContentType: "text/plain",
											Body:   msg     ,
										})
									failOnError(err, "Failed to publish a message")
								}else{
									fmt.Println("Cached")

								}
							}
						}

					}
				}

			}()


		}
	}


}
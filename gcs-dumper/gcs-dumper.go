package main
import (
	"os"
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"io/ioutil"
	"archive/zip"
	"path/filepath"
	"strings"
	"io"
	"time"
	"strconv"
	storage "google.golang.org/api/storage/v1"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"github.com/satori/go.uuid"

)


const (
	MAX_BYTES = 1024  * 1024 * 100
	scope = storage.DevstorageFullControlScope
	SECONDS = 30 * 60
)


var (
	bucketName  = os.Getenv("BUCKET")
	exchange = os.Getenv("EXCHANGE")

)



/** STOLEN FROM STACK OVERFLOW, PLEASE OPTIMIZE SHIT.
 */
func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}

}

func dumpFile(msgs [][]byte )  {
	for i := 0; i < len(msgs); i++ {
		err:=ioutil.WriteFile("tmp/"+string(uuid.NewV4()[15]), msgs[i], 0644)
		fmt.Println(err)
		//		w.Close()

	}
	t:=time.Now().Unix()
	zipit("tmp/","upload/data-"+strconv.FormatInt(t,10)+".zip")
	RemoveContents("tmp")
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func upload(){
	fmt.Println("upload time")
	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}
	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}
	filepath.Walk("upload", func(path string, info os.FileInfo, err error) error {
		if path == "upload"{
			return nil
		}
		// Insert an object into a bucket.
		fmt.Println(info.Name())
		object := &storage.Object{Name: info.Name()}
		file, _ := os.Open(path)
		if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
			fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			log.Fatal(err)
		}
		return err
	})
	RemoveContents("upload")


}



func main() {
	if bucketName == "" || exchange == ""{
		panic("Properties not set")
	}
	err:=os.Mkdir("tmp",0777)
	err=os.Mkdir("upload",0777)
	fmt.Println(err)
	RemoveContents("tmp")


	var total int = 0
	var data [][]byte
	compress:=make(chan [][]byte)
	uploadTime := time.NewTicker(time.Second * SECONDS).C

	//set := make(map[string]bool)
	rabbit:=os.Getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR")
	if rabbit == ""{
		panic("No environment for rabbit")
	}
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
		exchange, // exchange
		false,
		nil)

	failOnError(err, "Failed to bind to the queue.")

	msgs, err := ch.Consume(
		q.Name, // queue
		"s3-dumber",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			size:=len(d.Body)
			total+=size
			data = append(data, d.Body)
			fmt.Println(total)
			if err != nil {
				fmt.Println("Failed to parse tweet")
			}
			if total >= MAX_BYTES {
				compress<-data
				data=nil
				total=0
				failOnError(err,"ouch")
			}

		}
	}()

	for {
		select {
		case t:=<-compress:
		go func(){
			dumpFile(t)
		}()
		case <-uploadTime:
		go func(){
			upload()
		}()

		}
	}

}




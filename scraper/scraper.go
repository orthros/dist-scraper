package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

const (
	imgSelecter   = "#img"
	totalSelecter = "#selectpage"
	nextSelecter  = "#imgholder > a:nth-child(1)"
	totalRegex    = `of (\d+)`
)

func scrape(baseURL string, bookNameChapter string, hook FoundImageHook) {
	r := regexp.MustCompile(totalRegex)

	totalPages := 1
	nextLocation := baseURL + bookNameChapter
	for i := 0; i < totalPages; i++ {
		doc, err := goquery.NewDocument(nextLocation)
		failOnError(err, "Could not navigate")

		matches := r.FindAllString(doc.Find(totalSelecter).First().Text(), -1)
		totalPages, err = strconv.Atoi(matches[1])
		failOnError(err, "Could not determine the total pages")

		//Now to the image
		imagePath, _ := doc.Find(imgSelecter).Attr("src")

		data := getImageData(imagePath)

		hook.found(i, data)

		nextPath, _ := doc.Find(nextSelecter).Attr("href")
		nextLocation = baseURL + nextPath
	}
}

func getImageData(imagePath string) []byte {
	resp, err := http.Get(imagePath)
	failOnError(err, "Could not download the image")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Could not read the response")

	return body
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rabbitmqServer := os.Getenv("QUEUE_LOCATION")
	queueName := os.Getenv("QUEUE_NAME")

	conn, err := amqp.Dial(rabbitmqServer)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var message Message
			if err := json.Unmarshal(d.Body, &message); err != nil {
				panic(err)
			}

			//Todo: make this the non void one
			foundImageHook := &VoidFoundImageHook{}

			//Combine the two to get viable starting URL
			bookNameChapter := message.BookName + "/" + string(message.ChapterNumber)

			//Begin the scraping
			scrape(message.BaseUrl, bookNameChapter, foundImageHook)

			//Done scraping, log and ack
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

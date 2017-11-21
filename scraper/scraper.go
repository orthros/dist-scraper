package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func scrape(baseUrl string, bookNameChapter string, hook FoundImageHook) {
	const imgSelecter = "#img"
	const totalSelecter = "#selectpage"
	const nextSelecter = "#imgholder > a:nth-child(1)"

	totalPages := 1
	nextLocation := baseUrl + bookNameChapter
	for i := 0; i < totalPages; i++ {
		doc, err := goquery.NewDocument(nextLocation)
		failOnError(err, "Could not navigate")

		doc.Find(totalSelecter).First().Each(func(j int, s *goquery.Selection) {
			//Set the Total pages
		})

		//Now to the image
		imagePath, _ := doc.Find(imgSelecter).Attr("src")

		data := getImageData(imagePath)

		hook.found(i, data)

		nextPath, _ := doc.Find(nextSelecter).Attr("href")
		nextLocation = baseUrl + nextPath
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

type FoundImageHook interface {
	found(pageNum int, data []byte)
}

type ServiceFoundImageHook struct {
	ChapterID int
}

func (sfih ServiceFoundImageHook) found(pageNum int, data []byte) {
	postImage(sfih.ChapterID, pageNum, data)
}

type VoidFoundImageHook struct {
}

func (vfih VoidFoundImageHook) found(pageNum int, data []byte) {
	log.Printf("Found an image")
}

type Book struct {
	ID    int
	Title string
}

type Chapter struct {
	ID            int
	BookID        int
	ChapterNumber int
	ChapterTitle  string
}

type Page struct {
	ID         int
	ChapterID  int
	PageNumber int
	Data       []byte
}

type Message struct {
	BaseUrl       string
	BookName      string
	ChapterNumber int
}

func getBookID(bookName string) int {
	endpoint := os.Getenv("BOOKS_ENDPOINT")

	targetBook := &Book{
		ID:    0,
		Title: bookName,
	}
	jsonPage, err := json.Marshal(targetBook)
	failOnError(err, "Unable to marshal Book to json")

	resp, err := http.Post(endpoint+"/books", "applicaiton/json", bytes.NewBuffer(jsonPage))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	data := binary.BigEndian.Uint64(body)
	return int(data)
}

func getChapterID(bookID int, chapterNumber int) int {
	endpoint := os.Getenv("BOOKS_ENDPOINT")

	targetChapter := &Chapter{
		ID:            0,
		ChapterTitle:  "",
		BookID:        bookID,
		ChapterNumber: chapterNumber,
	}
	jsonBook, err := json.Marshal(targetChapter)
	failOnError(err, "Unable to marshal Chapter to json")

	resp, err := http.Post(endpoint+"/chapters", "applicaiton/json", bytes.NewBuffer(jsonBook))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	data := binary.BigEndian.Uint64(body)
	return int(data)
}

func postImage(chapterID int, pageNumber int, pageData []byte) {
	endpoint := os.Getenv("BOOKS_ENDPOINT")

	targetPage := &Page{
		ID:         0,
		ChapterID:  chapterID,
		PageNumber: pageNumber,
		Data:       pageData,
	}
	jsonPage, err := json.Marshal(targetPage)
	failOnError(err, "Unable to marshal Book to json")

	resp, err := http.Post(endpoint+"/pages", "applicaiton/json", bytes.NewBuffer(jsonPage))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	log.Printf("Got some data %s", body)
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

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

type BookService struct {
	endpoint string
}

func NewBookService() BookService {
	endpoint := os.Getenv("BOOKS_ENDPOINT")
	return BookService{
		endpoint: endpoint,
	}
}

func (service BookService) getBookID(bookName string) int {
	targetBook := &Book{
		ID:    0,
		Title: bookName,
	}
	jsonPage, err := json.Marshal(targetBook)
	failOnError(err, "Unable to marshal Book to json")

	resp, err := http.Post(service.endpoint+"/books", "applicaiton/json", bytes.NewBuffer(jsonPage))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	data := binary.BigEndian.Uint64(body)
	return int(data)
}

func (service BookService) getChapterID(bookID int, chapterNumber int) int {
	targetChapter := &Chapter{
		ID:            0,
		ChapterTitle:  "",
		BookID:        bookID,
		ChapterNumber: chapterNumber,
	}
	jsonBook, err := json.Marshal(targetChapter)
	failOnError(err, "Unable to marshal Chapter to json")

	resp, err := http.Post(service.endpoint+"/chapters", "applicaiton/json", bytes.NewBuffer(jsonBook))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	data := binary.BigEndian.Uint64(body)
	return int(data)
}

func (service BookService) postImage(chapterID int, pageNumber int, pageData []byte) {
	targetPage := &Page{
		ID:         0,
		ChapterID:  chapterID,
		PageNumber: pageNumber,
		Data:       pageData,
	}
	jsonPage, err := json.Marshal(targetPage)
	failOnError(err, "Unable to marshal Book to json")

	resp, err := http.Post(service.endpoint+"/pages", "applicaiton/json", bytes.NewBuffer(jsonPage))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	log.Printf("Got some data %s", body)
}

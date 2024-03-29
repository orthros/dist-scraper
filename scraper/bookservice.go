package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	endpointKey = "BOOKS_ENDPOINT"
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

type bookService struct {
	endpoint string
}

func newBookService() bookService {
	endpoint := os.Getenv(endpointKey)
	return bookService{
		endpoint: endpoint,
	}
}

func (service bookService) getBookID(bookName string) int {
	targetBook := &Book{
		ID:    0,
		Title: bookName,
	}
	jsonPage, err := json.Marshal(targetBook)
	failOnError(err, "Unable to marshal Book to json")

	resp, err := http.Post(service.endpoint+"/books", "application/json", bytes.NewBuffer(jsonPage))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	var dta float64
	err = json.Unmarshal(body, &dta)

	return int(dta)
}

func (service bookService) getChapterID(bookID int, chapterNumber int) int {
	targetChapter := &Chapter{
		ID:            0,
		ChapterTitle:  "",
		BookID:        bookID,
		ChapterNumber: chapterNumber,
	}
	jsonBook, err := json.Marshal(targetChapter)
	failOnError(err, "Unable to marshal Chapter to json")

	resp, err := http.Post(service.endpoint+"/chapters", "application/json", bytes.NewBuffer(jsonBook))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	var dta float64
	err = json.Unmarshal(body, &dta)

	return int(dta)
}

func (service bookService) postImage(chapterID int, pageNumber int, pageData []byte) {
	targetPage := &Page{
		ID:         0,
		ChapterID:  chapterID,
		PageNumber: pageNumber,
		Data:       pageData,
	}
	jsonPage, err := json.Marshal(targetPage)
	failOnError(err, "Unable to marshal Book to json")

	resp, err := http.Post(service.endpoint+"/pages", "application/json", bytes.NewBuffer(jsonPage))
	failOnError(err, "Couldn't get to endpoint")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Error reading response body")

	log.Printf("Posted image and got response: %s", body)
}

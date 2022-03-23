package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func getUrl(result chan string, url string) {
	response, errResponse := http.Get(url)

	if errResponse != nil {
		log.Println(errResponse)
		return
	}

	fmt.Println(response)
	urlContent, errUrl := io.ReadAll(response.Body)
	if errUrl != nil {
		log.Println(errUrl)
		return
	}

	result <- string(urlContent)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)
}

func main() {
	f, err := os.OpenFile("save.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	const BaseUrl = "https://reqres.in/api"
	urls := []string{BaseUrl + "/user?page=2", BaseUrl + "/user/2", BaseUrl + "/unknown"}

	urlChannel := make(chan string)
	urlChannel2 := make(chan string)
	urlChannel3 := make(chan string)
	urlChannels := []chan string{urlChannel, urlChannel2, urlChannel3}

	for index, url := range urls {
		go getUrl(urlChannels[index], url)
	}

	for _, urlChannel := range urlChannels {
		if _, err := f.WriteString(<-urlChannel); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return
			}
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

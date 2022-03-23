package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func countWordInFile(result chan int, path string, word string) {
	count := 0
	content, err := os.ReadFile(path)

	if err != nil {
		result <- 0
		log.Fatalln(err)
	}

	count = strings.Count(string(content), word)
	fmt.Println(count)
	result <- count
	defer close(result)
}

func main() {
	wordCounterChannel := make(chan int)
	wordCounterChannel2 := make(chan int)

	go countWordInFile(wordCounterChannel, "text.txt", "Ipsum")
	go countWordInFile(wordCounterChannel2, "text2.txt", "Ipsum")

	fmt.Println(<-wordCounterChannel + <-wordCounterChannel2)
}

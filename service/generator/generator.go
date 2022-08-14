package generator

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

const ADJ_LENGTH int = 1347
const NOUN_LENGTH int = 1525

var adjectives []string = make([]string, ADJ_LENGTH)
var nouns []string = make([]string, NOUN_LENGTH)

// Loads a provided slice with strings from a file referenced by the filepath
func loadSlice(filepath string, slice []string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var counter int = 0
	for scanner.Scan() {
		slice[counter] = scanner.Text()
		counter++
	}
	file.Close()

	return nil
}

// Asynchronously runs loadSlice functions to populate approperiate word slices
func LoadWords() error {
	aCh := make(chan error)
	nCh := make(chan error)

	go func() {
		aCh <- loadSlice("./service/generator/english-adjectives.txt", adjectives)
	}()
	go func() {
		nCh <- loadSlice("./service/generator/english-nouns.txt", nouns)
	}()

	if err := <-aCh; err != nil {
		return err
	}
	if err := <-nCh; err != nil {
		return err
	}
	return nil
}

// Returns a randomly generated adjective-word string from provided word slices
func GenerateURL() string {
	rand.Seed(time.Now().UnixNano())
	adjIndex := rand.Intn(ADJ_LENGTH)
	nounIndex := rand.Intn(NOUN_LENGTH)
	return adjectives[adjIndex] + "-" + nouns[nounIndex]
}

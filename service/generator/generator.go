package generator

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

const ADJ_MAX_LENGTH int = 1313
const NOUN_MAX_LENGTH int = 1525

var adj_count int = ADJ_MAX_LENGTH
var noun_count int = NOUN_MAX_LENGTH

var adjectives []string = make([]string, ADJ_MAX_LENGTH)
var nouns []string = make([]string, NOUN_MAX_LENGTH)

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
	rand.Seed(time.Now().UnixNano())

	aCh := make(chan error)
	nCh := make(chan error)

	go func() {
		aCh <- loadSlice("./dict/english-adjectives.txt", adjectives)
	}()
	go func() {
		nCh <- loadSlice("./dict/english-nouns.txt", nouns)
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
	adjIndex := rand.Intn(ADJ_MAX_LENGTH)
	nounIndex := rand.Intn(NOUN_MAX_LENGTH)
	var result string = adjectives[adjIndex] + "_" + nouns[nounIndex]
	removeWord(&adjectives, adjIndex, &adj_count)
	removeWord(&nouns, nounIndex, &noun_count)
	return result
}

func BreakURL(memorableUrl string) {
	words := strings.Split(memorableUrl, "_")
	appendWord(&adjectives, words[0], &adj_count)
	appendWord(&nouns, words[1], &noun_count)
}

func removeWord(slice *[]string, index int, counter *int) {
	(*slice)[index] = (*slice)[len(*slice)-1]
	*counter--
	*slice = (*slice)[:len(*slice)-1]
}

func appendWord(slice *[]string, word string, counter *int) {
	*slice = append(*slice, word)
	*counter++
}

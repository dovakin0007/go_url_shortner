package routes

import (
	"log"
	"math/rand"
	"time"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCEFGHIJKLMNOPQRSTUVW1234567890")

func randSeq(n int) string {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[newRand.Intn(len(letters))]
	}
	return string(b)
}

func err_handler(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

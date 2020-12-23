package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func WriteFileLog(whichfile string, s string) {
	t := time.Now().Format(time.UnixDate)
	title := fmt.Sprintf("%v : %v\n", t, s)
	f, err := os.OpenFile(whichfile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(title); err != nil {
		log.Println(err)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomString(n int32) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

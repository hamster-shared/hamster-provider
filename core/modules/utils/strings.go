package utils

import (
	"math/rand"
	"strings"
	"time"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Remove(slice []string, e string) []string {
	index := IndexOf(slice, e)
	if index >= 0 {
		return RemoveIndex(slice, index)
	}
	return slice
}

func RemoveIndex(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func IndexOf(slice []string, e string) int {
	for i, a := range slice {
		if a == e {
			return i
		}
	}
	return -1
}

func GetRandomString(length int) string {
	rand.Seed(time.Now().Unix())

	//Only lowercase
	charSet := "abcdedfghijklmnopqrst"
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

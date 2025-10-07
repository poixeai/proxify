package util

import (
	"math/rand"
	"time"
)

var RandLowercaseAlphanumeric = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
var RandAlphanumeric = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var RandNumeric = []rune("0123456789")

// RandGenerater generates a random string of length n from the given runes
func RandGenerater(runes []rune, n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]rune, n)
	for i := range b {
		b[i] = runes[r.Intn(len(runes))]
	}
	return string(b)
}

func GenerateRequestID() string {
	now := time.Now().Format("20060102150405")
	randomPart := RandGenerater(RandNumeric, 6)

	return now + randomPart
}

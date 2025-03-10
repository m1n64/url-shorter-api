package utils

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func PtrToTimeString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	str := t.String()
	return &str
}

func RandStringBytesRmndr(n int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[random.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

package gologger

import (
	"crypto/rand"
	"fmt"
	"time"
)

func generateSafeToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	byteSlice := make([]byte, length)
	_, err := rand.Read(byteSlice)
	if err != nil {
		return ""
	}
	for index, the_byte := range byteSlice {
		byteSlice[index] = charset[the_byte%byte(len(charset))]
	}
	return string(byteSlice)
}

func generateFileID() string {
	safeToken := generateSafeToken(5)
	timestamp := time.Now().Unix()
	ID := fmt.Sprintf("%d-%s", timestamp, safeToken)
	return ID
}

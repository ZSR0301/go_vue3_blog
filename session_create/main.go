package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomKey(length int) (string, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

func main() {
	// 生成一个长度为 32 字节的随机密钥
	keyLength := 32
	randomKey, err := GenerateRandomKey(keyLength)
	if err != nil {
		panic(err)
	}
	fmt.Println("Randomly generated key:", randomKey)
}

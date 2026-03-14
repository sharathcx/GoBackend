package utils

import (
	"math/rand/v2"
)


func InvokeUID(prefix string, length int) string {
	if prefix == "" {
		prefix = "uid"
	}
	if length <= 0 {
		length = 10
	}

	charPool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	uid := make([]byte, length)

	for i := 0; i < length; i++ {
		uid[i] = charPool[rand.IntN(len(charPool))]
	}

	return prefix + "_" + string(uid)
}

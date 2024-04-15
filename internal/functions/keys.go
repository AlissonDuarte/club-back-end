package functions

import (
	"math/rand"
	"time"
)

func GenerateKeys(n int) string {
	var letras = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letras[rand.Intn(len(letras))]
	}
	return string(b)
}

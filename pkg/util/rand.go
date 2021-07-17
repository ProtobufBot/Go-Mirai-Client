package util

import (
	"math/rand"
	"strings"
)

func RandomString(randGen *rand.Rand,len int) string {
	return RandomStringRange(randGen,len, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func RandomStringRange(randGen *rand.Rand,length int, str string) string {
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteByte(str[randGen.Intn(len(str))])
	}
	return sb.String()
}

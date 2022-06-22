package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}

func RandomNumber() int64 {
	return RandomInt(0, 1000000)
}

//func RandomGender() string {
//	gender := []string{"male", "female"}
//	n := len(gender)
//	return gender[rand.Intn(n)]
//}

func RandomRole() string {
	role := []string{"admin", "partner", "locket"}
	n := len(role)
	return role[rand.Intn(n)]
}

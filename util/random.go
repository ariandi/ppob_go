package util

import (
	"fmt"
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

func RandomNumberStr() string {
	rInt := RandomInt(1, 999999)
	s := fmt.Sprintf("%010d", rInt)
	return s
}

//func RandomGender() string {
//	gender := []string{"male", "female"}
//	n := len(gender)
//	return gender[rand.Intn(n)]
//}

func RandomRole() string {
	role := []string{"Admin", "Biller", "Merchant"}
	n := len(role)
	return role[rand.Intn(n)]
}

func RandomStatus() string {
	status := []string{"active", "inactive", "trouble"}
	n := len(status)
	return status[rand.Intn(n)]
}

func RandomTrxStatus() string {
	status := []string{"success", "pending", "failed", "settlement", "reversal"}
	n := len(status)
	return status[rand.Intn(n)]
}

func SetTxID() string {
	txID := time.Time.Format(time.Now(), "200601021504") + padLeft(RandomNumberUnique())
	return txID
}

func padLeft(nr int64) string {
	s := fmt.Sprintf("%04d", nr)
	return s
}

func RandomNumberUnique() int64 {
	return RandomInt(1, 9999)
}

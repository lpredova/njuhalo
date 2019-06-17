package helper

import (
	"fmt"
	"math/rand"
	"time"
)

const dateLayout = "02.01.2006 15:04:05"

// TimestampToDate converts unit timestamp to Croatian date format
func TimestampToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.UTC().Format(dateLayout)
}

// RandomString generates random string
func RandomString() string {
	n := rand.Intn(20)
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

package helper

import (
	"fmt"
	"math/rand"
	"regexp"
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

// GetSliceData returns data specified by index if exists
func GetSliceData(slice []string, index int) string {
	if len(slice) > (index + 1) {
		return slice[index]
	}

	return ""
}

// GetNumber returns numeric value from strong if there is any
func GetNumber(input string) string {
	re := regexp.MustCompile("[0-9]+")
	result := re.FindAllString(input, -1)
	if len(result) > 0 {
		return result[0]
	}

	return ""
}

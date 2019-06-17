package helper

import (
	"time"
)

const dateLayout = "02.01.2006 15:04:05"

// TimestampToDate converts unit timestamp to Croatian date format
func TimestampToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.UTC().Format(dateLayout)
}

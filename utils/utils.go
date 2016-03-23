package utils

import "github.com/satori/go.uuid"
import "time"

func Epoch_ms() int64 {
	now := time.Now()
	return now.UnixNano() / 1000000
}

func Makeid() string {
	return uuid.NewV4().String()
}

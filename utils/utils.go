package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/satori/go.uuid"
	"time"
)

func Epoch_ms() int64 {
	now := time.Now()
	return now.UnixNano() / 1000000
}

func Makeid() string {
	return uuid.NewV4().String()
}

func Sign(s ...string) string {
	mac := hmac.New(sha256.New, []byte("thisshouldbeasecretkey"))

	for i := 0; i < len(s); i++ {
		mac.Write([]byte(s[i]))
	}

	sum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

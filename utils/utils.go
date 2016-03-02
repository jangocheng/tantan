package utils

import (
	math_rand "math/rand"
	"time"
)

func init() {
	math_rand.Seed(time.Now().Unix())
}

func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func GetNowSecond() int {
	return int(time.Now().Unix())
}

func GetNowStringYMD() string {
	return time.Now().Format("2006-01-02")
}

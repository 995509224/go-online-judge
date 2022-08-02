package help

import (
	"math/rand"
	"strconv"
	"time"
)

func Getyzm() string {
	var str string
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 6; i++ {
		t := rand.Intn(10)
		str += strconv.Itoa(t)
	}
	return str
}

package lokahiadminserver

import (
	"log"
	"math/rand"
	"time"
)

func init () {
	rand.Seed(time.Now().UnixNano())

	dur := time.Duration(rand.Intn(15)) * time.Second
	log.Printf("sleeping for %v", dur)
	time.Sleep(dur)
}

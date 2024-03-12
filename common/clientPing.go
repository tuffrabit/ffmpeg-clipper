package common

import (
	"log"
	"time"
)

var clientLastUpdateTime time.Time

func SetClientLastUpdateTime(t time.Time) {
	clientLastUpdateTime = t
}

func MonitorClient() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		//fmt.Println("Tick at", t)
		now := time.Now()
		if now.Sub(clientLastUpdateTime).Seconds() > 10 {
			log.Fatal("Client no longer active.")
		}
	}
}

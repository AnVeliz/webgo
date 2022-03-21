package app

import (
	"fmt"
	"time"
)

func RunTimerAsync() <-chan time.Time {
	timeNowChan := make(chan time.Time)

	go func() {
		timer := time.NewTicker(time.Duration(1 * time.Second))
		for {
			timeNow := <-timer.C
			fmt.Printf("Timer fired: %s\n", timeNow)
			timeNowChan <- timeNow
		}
	}()

	return timeNowChan
}

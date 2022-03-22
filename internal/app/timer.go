package app

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func RunTimerAsync(ctx context.Context, wg *sync.WaitGroup) <-chan time.Time {
	timeNowChan := make(chan time.Time)

	wg.Add(1)
	go func() {
		timer := time.NewTicker(time.Duration(1 * time.Second))
		fmt.Println("Timer is running")
	loop:
		for {
			select {
			case timeNow := <-timer.C:
				fmt.Printf("Timer fired: %s\n", timeNow)
				timeNowChan <- timeNow
			case <-ctx.Done():
				break loop
			}
		}
		timer.Stop()
		fmt.Println("Timer has finished")
		wg.Done()
	}()

	return timeNowChan
}

//go:build windows

package app

import (
	"context"
	"fmt"
	"sync"
	"syscall"
	"time"
)

func StayAwakeAsync(ctx context.Context, wg *sync.WaitGroup) {
	var pulseTime = 10 * time.Second
	go stayAwake(ctx, wg, pulseTime)
}

func stayAwake(ctx context.Context, wg *sync.WaitGroup, pulseTime time.Duration) {
	const (
		esSystemRequired = 0x00000001
		esContinuous     = 0x80000000
	)

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setThreadExecStateProc := kernel32.NewProc("SetThreadExecutionState")

	pulse := time.NewTicker(pulseTime)

	wg.Add(1)
	fmt.Println("stay awake keeper is running")
loop:
	for {
		select {
		case <-pulse.C:
			fmt.Println("Stay awake pulse")
			setThreadExecStateProc.Call(uintptr(esSystemRequired))
		case <-ctx.Done():
			break loop
		}
	}
	fmt.Println("stay awake keeper has finished")
	pulse.Stop()
	wg.Done()
}

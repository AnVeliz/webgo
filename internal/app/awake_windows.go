//go:build windows

package app

import (
	"fmt"
	"syscall"
	"time"
)

func StayAwakeAsync() {
	var pulseTime = 10 * time.Second
	go stayAwake(pulseTime)
}

func stayAwake(pulseTime time.Duration) {
	const (
		esSystemRequired = 0x00000001
		esContinuous     = 0x80000000
	)

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setThreadExecStateProc := kernel32.NewProc("SetThreadExecutionState")

	pulse := time.NewTicker(pulseTime)

	for {
		select {
		case <-pulse.C:
			fmt.Println("Stay awake pulse")
			setThreadExecStateProc.Call(uintptr(esSystemRequired))
		}
	}
}

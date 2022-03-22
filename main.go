package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AnVeliz/webgo/internal/app"
	"github.com/AnVeliz/webgo/internal/chromium"
	"github.com/AnVeliz/webgo/internal/httpsrv"
	"github.com/AnVeliz/webgo/internal/wssrv"
)

var (
	//go:embed assets/*
	embeddedContent embed.FS
)

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	app.StayAwakeAsync(ctx, &wg)

	baseUrl := "http://localhost"
	httpPort := 8088
	wsPort := 8089

	mainWsSrv, writeToSocket, _ := runWebSocketAsync(ctx, &wg, wsPort)
	runMainAppAsync(ctx, &wg, writeToSocket)
	mainHttpSrv := runHttpServerAsync(&wg, httpPort)

	runChromiumAsync(ctx, cancel, &wg, baseUrl, httpPort, mainHttpSrv, mainWsSrv)

	fmt.Println("the app is running")
	wg.Wait()
	fmt.Println("the app has finished")
}

func runHttpServerAsync(wg *sync.WaitGroup, httpPort int) *http.Server {
	mainHttpSrv := httpsrv.Create(httpPort, embeddedContent)
	wg.Add(1)
	go func() {
		httpsrv.Run(mainHttpSrv)
		wg.Done()
	}()
	return mainHttpSrv
}

func runChromiumAsync(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, baseUrl string, httpPort int, mainHttpSrv *http.Server, mainWsSrv *http.Server) {
	wg.Add(1)
	go func() {
		chromium.Run(fmt.Sprintf("%s:%d/", baseUrl, httpPort))

		mainHttpSrv.Shutdown(ctx)
		mainWsSrv.Shutdown(ctx)
		cancel()
		wg.Done()
	}()
}

func runMainAppAsync(ctx context.Context, wg *sync.WaitGroup, writeToSocket chan<- []byte) {
	wg.Add(1)
	go func() {
		fmt.Println("main app async loop is running")
		timerChan := app.RunTimerAsync(ctx, wg)
	loop:
		for {
			select {
			case currentTime := <-timerChan:
				writeToSocket <- []byte(currentTime.Format(time.RFC3339))
			case <-ctx.Done():
				break loop
			}
		}
		fmt.Println("main app async loop has finished")
		wg.Done()
	}()
}

func runWebSocketAsync(ctx context.Context, wg *sync.WaitGroup, wsPort int) (*http.Server, chan<- []byte, <-chan []byte) {
	mainWsSrv, writeToSocket, readFromSocket := wssrv.Create(ctx, wg, wsPort)
	wg.Add(1)
	go func() {
		wssrv.Run(mainWsSrv)
		wg.Done()
	}()
	return mainWsSrv, writeToSocket, readFromSocket
}

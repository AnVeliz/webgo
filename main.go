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
	embeddedContentFS embed.FS
)

func main() {
	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.StayAwakeAsync(ctx, &wg)

	baseUrl := "http://localhost"
	httpPort := 8088
	wsPort := 8089

	mainWsSrv, writeToSocket, _ := runWebSocketAsync(ctx, &wg, wsPort)
	runMainAppAsync(ctx, &wg, writeToSocket)
	mainHttpSrv := runHttpServerAsync(&wg, httpPort)

	runChromiumAsync(ctx, stop, &wg, baseUrl, httpPort, mainHttpSrv, mainWsSrv)

	fmt.Println("the app is running")
	wg.Wait()
	fmt.Println("the app has finished")
}

func runHttpServerAsync(wg *sync.WaitGroup, httpPort int) *http.Server {
	mainHttpSrv := httpsrv.Create(httpPort, embeddedContentFS)
	wg.Add(1)
	go func() {
		httpsrv.Run(mainHttpSrv)
		wg.Done()
	}()
	return mainHttpSrv
}

func runChromiumAsync(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup, baseUrl string, httpPort int, mainHttpSrv *http.Server, mainWsSrv *http.Server) {
	wg.Add(1)
	go func() {
		chromium.Run(fmt.Sprintf("%s:%d/", baseUrl, httpPort), embeddedContentFS)

		mainHttpSrv.Shutdown(ctx)
		mainWsSrv.Shutdown(ctx)
		stop()
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

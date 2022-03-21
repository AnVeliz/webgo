package main

import (
	"context"
	"embed"
	"fmt"

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
	baseUrl := "http://localhost"
	httpPort := 8088
	wsPort := 8089

	mainWsSrv, writeToSocket, _ := wssrv.Create(wsPort)
	go func() {
		wssrv.Run(mainWsSrv)
	}()
	go func() {
		timerChan := app.RunTimerAsync()
		for {
			currentTime := <-timerChan
			writeToSocket <- []byte(currentTime.String())
		}
	}()

	mainHttpSrv := httpsrv.Create(httpPort, embeddedContent)

	go func() {
		chromium.Run(fmt.Sprintf("%s:%d/", baseUrl, httpPort))

		ctx := context.Background()
		mainHttpSrv.Shutdown(ctx)
		mainWsSrv.Shutdown(ctx)
	}()

	httpsrv.Run(mainHttpSrv)
}

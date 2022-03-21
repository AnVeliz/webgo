package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/AnVeliz/webgo/internal/chromium"
	"github.com/AnVeliz/webgo/internal/httpsrv"
)

var (
	//go:embed assets/*
	embeddedContent embed.FS
)

func main() {
	baseUrl := "http://localhost"
	port := 8088
	address := fmt.Sprintf("%s:%d/", baseUrl, port)

	srv := httpsrv.Create(port, embeddedContent)

	go func() {
		chromium.Run(address)
		srv.Shutdown(context.Background())
	}()

	httpsrv.Run(srv)
}

package wssrv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func Create(ctx context.Context, wg *sync.WaitGroup, port int) (*http.Server, chan<- []byte, <-chan []byte) {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	writeToSocket := make(chan []byte)
	readFromSocket := make(chan []byte)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Printf("can not upgrade HTTP to WebSocket: %s\n", err)
			return
		}

		wg.Add(1)
		go func() {
			defer conn.Close()
			fmt.Println("user's channel holder is running")

			var lwg sync.WaitGroup
			lwg.Add(2)
			go func() {
				fmt.Println("channel writer is running")
			loop:
				for {
					select {
					case msg, ok := <-writeToSocket:
						if !ok {
							fmt.Println("can not read from closed channel")
							break loop
						}
						if err := wsutil.WriteServerMessage(conn, ws.OpText, msg); err != nil {
							if _, ok := err.(wsutil.ClosedError); !ok {
								fmt.Printf("error: %s\n", err)
							}
							break loop
						}
						fmt.Printf("message sent: %s\n", msg)
					case <-ctx.Done():
						break loop
					}
				}
				fmt.Println("channel writer has finished")
				lwg.Done()
			}()
			go func() {
				fmt.Println("channel reader is running")
			loop:
				for {
					conn.SetReadDeadline(time.Now().Add(5 * time.Second))
					data, _, err := wsutil.ReadClientData(conn)
					if err != nil && !os.IsTimeout(err) {
						if _, ok := err.(wsutil.ClosedError); !ok {
							fmt.Printf("error: %s\n", err)
						}
						break loop
					} else if os.IsTimeout(err) {
						break loop
					}
					readFromSocket <- data
				}
				fmt.Println("channel reader has finished")
				lwg.Done()
			}()

			lwg.Wait()
			fmt.Println("user's channel holder has finished")
			wg.Done()
		}()
	}))
	log.Println("main WebSocket server handler has set up")

	return srv, writeToSocket, readFromSocket
}

func Run(srv *http.Server) {
	fmt.Println("main WebSocket server handler is running")
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Println("main WebSocket server handler has stopped")
}

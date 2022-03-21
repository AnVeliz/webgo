package wssrv

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func Create(port int) (*http.Server, chan<- []byte, <-chan []byte) {
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

		go func() {
			defer conn.Close()

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				for {
					msg, ok := <-writeToSocket
					if !ok {
						fmt.Println("can not read from closed channel")
						break
					}
					if err := wsutil.WriteServerMessage(conn, ws.OpText, msg); err != nil {
						fmt.Printf("error: %s\n", err)
						break
					}
					fmt.Printf("message sent: %s\n", msg)
				}
				wg.Done()
			}()
			go func() {
				for {
					data, _, err := wsutil.ReadClientData(conn)
					if err != nil {
						fmt.Printf("error: %s\n", err)
						break
					}
					readFromSocket <- data
				}
				wg.Done()
			}()

			wg.Wait()
		}()
	}))
	log.Println("main WebSocket server handler has set up")

	return srv, writeToSocket, readFromSocket
}

func Run(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Println("main WebSocket server handler has stopped")
}

package httpsrv

import (
	"embed"
	"fmt"
	"log"
	"net/http"
)

func Create(port int, fs embed.FS) *http.Server {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	mux.Handle("/", http.FileServer(http.FS(fs)))
	log.Println("main file server handler has set up")

	return srv
}

func Run(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Println("main file server handler has stopped")
}

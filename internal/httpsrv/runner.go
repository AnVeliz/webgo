package httpsrv

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

func Create(port int, eFS embed.FS) *http.Server {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	subEFS, err := fs.Sub(eFS, "assets/web-ui-public-build/public")
	if err != nil {
		fmt.Printf("can not open sub filesystem to serve webui, error: %s\n", err)
		return nil
	}

	mux.Handle("/", http.FileServer(http.FS(subEFS)))
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

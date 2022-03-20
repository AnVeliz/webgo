package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var (
	//go:embed webui/my-new-app/out/my-new-app-win32-x64/*
	res embed.FS
)

func main() {
	http.FileServer(http.FS(res))
	log.Println("server started...")

	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})*/
	http.Handle("/", http.FileServer(http.FS(res)))

	go func() {
		cmnd := exec.Command("http://localhost:8088/webui/my-new-app/out/my-new-app-win32-x64/my-new-app.exe", "")
		//cmnd.Run() // and wait
		err := cmnd.Start()
		if err != nil {
			fmt.Println(err)
		}
	}()

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}

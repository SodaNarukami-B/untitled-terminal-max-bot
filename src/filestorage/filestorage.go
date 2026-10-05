package filestorage

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var STORAGE string = "repo"

func Run() {

	http.HandleFunc("GET /repo/{token}", func(resw http.ResponseWriter, req *http.Request) {
		token := req.PathValue("token")

		log.Printf("<<< GET %v\n", token)

		saveFilename := filepath.Base(token)
		filePath := filepath.Join(STORAGE, saveFilename)

		info, err := os.Stat(filePath)

		if os.IsNotExist(err) || info.IsDir() {
			http.Error(resw, "File not found", http.StatusNotFound)
			return
		}

		http.ServeFile(resw, req, filePath)
	})

	go func() {
		log.Printf("[ FILESTORAGE STARTED ]\n\n")
		log.Fatalln(http.ListenAndServe(":60002", nil))
	}()

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)

	<-exit

	fmt.Printf("\r")
	log.Print("Exiting...")
}

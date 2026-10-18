package server

import (
	"encoding/json"
	"fmt"
	"log"

	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mime"
	"net/http"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

var MCU_CommandHandlers = map[string]MCU_COMMAND_HANDLER{
	"echo":   MCU_echo,
	"create": MCU_create,
}

func main() {
	fmt.Println("It's a module calling from master. Please use Master")
}

func Run() {

	dataBaseClient := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	http.HandleFunc("/webhook", func(resw http.ResponseWriter, req *http.Request) {
		// Http validation
		if req.Method != http.MethodPost {
			resw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if req.ContentLength <= 0 {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		var Update schemes.Update

		if json.Unmarshal(reqBody, &Update) != nil {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		switch Update.UpdateType {
		case schemes.TypeMessageCreated:
			var Message schemes.MessageCreatedUpdate

			if json.Unmarshal(reqBody, &Message) != nil {
				resw.WriteHeader(http.StatusBadRequest)
				return
			}

			text := Message.Message.Body.Text
			if len(text) == 0 {
				resw.WriteHeader(http.StatusBadRequest)
				return
			}

			command, err := ParseCommand(text)
			if err != nil {
				resw.WriteHeader(http.StatusBadRequest)
				resw.Write([]byte(err.Error()))
				return
			}

			handler, ex := MCU_CommandHandlers[command.Command]
			if !ex {
				resw.WriteHeader(http.StatusNotFound)
				resw.Write([]byte("Command not found or not implemented"))
				return
			}

			handler(resw, &Message.Message, *command, &dataBaseClient)

			return
		}
	})

	go func() {
		log.Println("Server runned")
		log.Fatalf("%v", http.ListenAndServe(":60001", nil))
	}()

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	sig := <-exit

	fmt.Printf("\r")
	log.Printf("Received %v. Exiting...\n", sig)
}

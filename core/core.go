package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Command struct {
	command  string
	args     []string
	raw_args string
}

type message_com_handler func(resw http.ResponseWriter, req *http.Request, msg *schemes.Message)

var MESSAGE_COMMAND = map[string]message_com_handler{
	"echo":   handle_echo,
	"create": handle_create,
}

func main() {
	http.HandleFunc("/webhook", func(resw http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			resw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req_body, err := io.ReadAll(req.Body)
		if err != nil {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		var upd schemes.Update

		if json.Unmarshal(req_body, &upd) != nil {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		switch upd.UpdateType {
		case schemes.TypeMessageCreated:
			{
				var msg schemes.MessageCreatedUpdate

				if json.Unmarshal(req_body, &msg) != nil {
					resw.WriteHeader(http.StatusBadRequest)
					return
				}

				log.Printf("< %s/n", msg.Message.Body.Text)

				command, err := parse_command(msg.Message.Body.Text)

				if err != nil {
					http.Error(resw, "Invalid command", http.StatusBadRequest)
					return
				}

				command_handler := MESSAGE_COMMAND[command.command]
				if command_handler == nil {
					http.Error(resw, "Not implemented command", http.StatusNotImplemented)

					return
				}

				command_handler(resw, req, &msg.Message)

				return
			}
		}
	})

	go func() {
		log.Fatal(http.ListenAndServe(":60001", nil))
	}()

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)

	<-exit

	fmt.Printf("\r")
	log.Print("Received exit signal. Exiting...")
}

func parse_command(data string) (*Command, error) {
	if len(data) == 0 {
		return nil, errors.New("No data")
	}

	if data[0] != '/' {
		return nil, errors.New("Not a command")
	}

	data = data[1:]

	_, raw_args, _ := strings.Cut(data, " ")

	parts := strings.Fields(data)

	return &Command{
		command:  parts[0],
		args:     parts[1:],
		raw_args: raw_args,
	}, nil
}

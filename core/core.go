package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

// TYPES
type Command struct {
	name     string
	args     []string
	raw_args string
}

// handling tables
type text_command_handler func(http.ResponseWriter, *http.Request, *schemes.Message)

var text_command_table = map[string]text_command_handler{
	"echo": handle_echo,
	// create soon
}

func main() {

	flag.Parse()

	text_command_table := map[string]text_command_handler{
		"echo": handle_echo,
	}

	http.HandleFunc("/webhook", func(resw http.ResponseWriter, req *http.Request) {
		// At this moment we can parse only post requests
		if req.Method != http.MethodPost {
			resw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Searching for method

		var msg schemes.MessageCreatedUpdate // from that type we need only UpdateType

		if json.NewDecoder(req.Body).Decode(&msg) != nil {
			resw.WriteHeader(http.StatusBadRequest)
			resw.Write([]byte("Invalid message"))
			return
		}

		// Parings received method
		switch msg.UpdateType {
		case schemes.TypeMessageCreated:
			{
				// MessageCreated can be used only for text command as analog for callback
				// Command parsing
				raw_text := msg.Message.Body.Text

				log.Printf("< %v\n", raw_text)

				comm, err := parse_command(raw_text)

				if err != nil {
					resw.WriteHeader(http.StatusBadRequest)
					resw.Write([]byte("Command expected"))
					return
				}

				// Handling command

				command_handler := text_command_table[comm.name]
				if command_handler == nil {
					resw.WriteHeader(http.StatusNotImplemented)
					resw.Write([]byte("Not implemented command"))
					return
				}

				command_handler(resw, req, &msg.Message)

				return
			}
		default:
			{
				resw.WriteHeader(http.StatusNotImplemented)
				resw.Write([]byte("Not implemented update type"))
				return
			}
		}
	})

	go func() {
		log.Fatal(http.ListenAndServe(":60001", nil))
	}()

	// Exit endpoint
	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	fmt.Print("\nExiting...\n")
}

func parse_command(text string) (*Command, error) {
	if len(text) == 0 {
		return nil, errors.New("Empty command")
	}

	if text[0] != '/' {
		return nil, errors.New("Not a command")
	}

	parts := strings.Fields(text)

	comm := strings.TrimPrefix(parts[0], "/")

	_, raw_args, _ := strings.Cut(text, " ")

	return &Command{
		name:     comm,
		args:     parts[1:],
		raw_args: raw_args,
	}, nil
}

package server

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
	Command  string
	Args     []string
	Raw_args string
}

type CommandHandler func(http.ResponseWriter, *http.Request, *schemes.Message, Command)

var MCU_commands = map[string]CommandHandler{
	"echo":   MCU_echo_handler,
	"create": MCU_create_handler,
}

func Run() {

	http.HandleFunc("/webhook", func(resw http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(resw, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		var upd schemes.Update

		if json.Unmarshal(reqBody, &upd) != nil {
			resw.WriteHeader(http.StatusBadRequest)
			return
		}

		switch upd.UpdateType {
		case schemes.TypeMessageCreated:
			{
				var msg schemes.MessageCreatedUpdate

				if json.Unmarshal(reqBody, &msg) != nil {
					resw.WriteHeader(http.StatusBadRequest)
					return
				}

				log.Printf("<<< %s\n", msg.Message.Body.Text)

				command, err := parse_command(msg.Message.Body.Text)

				if err != nil {
					http.Error(resw, err.Error(), http.StatusBadRequest)
					return
				}

				handler, exists := MCU_commands[command.Command]
				if !exists {
					http.Error(resw, "Command not found or not implemented", http.StatusNotImplemented)
					return
				}

				handler(resw, req, &msg.Message, *command)

				return
			}
		default:
			{
				resw.WriteHeader(http.StatusNotImplemented)
				return
			}
		}
	})

	go func() {
		log.Printf("[ SERVER STARTED ]\n")
		log.Fatalln(http.ListenAndServe(":60001", nil))
	}()

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	fmt.Printf("\r")
	log.Printf("Shutdowning...\n")
}

func parse_command(text string) (*Command, error) {
	if len(text) == 0 {
		return nil, errors.New("Empty command")
	}

	if text[0] != '/' {
		return nil, errors.New("Not a command. Commands MUST starting on '/'")
	}

	text = text[1:]

	command, raw_args, _ := strings.Cut(text, " ")
	args := strings.Fields(raw_args)

	return &Command{
		Command:  command,
		Args:     args,
		Raw_args: raw_args,
	}, nil
}

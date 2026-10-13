package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Command struct {
	Command  string   `json:"command"`
	Args     []string `json:"args"`
	Raw_args string   `json:"raw_args"`
}

var DefaultMessages = map[string]any{
	"message_created": schemes.MessageCreatedUpdate{
		Update: schemes.Update{
			UpdateType: "message_created",
			Timestamp:  int(time.Now().Unix()),
		},
		Message: schemes.Message{
			Sender: schemes.User{
				UserId:    17000000,
				Name:      "soda",
				Username:  "SodaNarukami",
				FirstName: "Soda",
				LastName:  "Narukami",
				IsBot:     false,
			},
			Recipient: schemes.Recipient{
				ChatId:   10000000,
				ChatType: schemes.DIALOG,
			},
			Timestamp: time.Now().Unix(),
			Body: schemes.MessageBody{
				Mid:  "apple-juice",
				Seq:  0,
				Text: "",
			},
		},
	},
}

var Message_client_commands = map[string]func(*schemes.Message, string) error{
	"create": COM_create_message,
}

func Run() {

	fmt.Print("[ CLIENT STARTED ]\n* If you see message starts with '::' - its a system messages\n\n")

	message_client()
}

func message_client() {
	client := http.Client{
		Timeout: time.Millisecond * 200,
	}

	msg, exists := DefaultMessages["message_created"].(schemes.MessageCreatedUpdate)
	if !exists {
		fmt.Printf(":: Cannot find message_created signature\n")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")

		if !scanner.Scan() {
			break
		}

		text := scanner.Text()

		command, err := parse_command(text)

		if err != nil {
			fmt.Printf(":: Failed to parse command - %v\n", err)
			continue
		}

		if command.Command == "quit" || command.Command == "exit" {
			break
		}

		handler, exists := Message_client_commands[command.Command]

		if exists {
			err := handler(&msg.Message, command.Raw_args)
			if err != nil {
				fmt.Printf(":: Failed to handle supported command '%s', sending as text ( %v )\n", command.Command, err)
				msg.Message.Body.Text = text
			}
		} else {
			msg.Message.Body.Text = text
		}

		json_data, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf(":: Marshaling final paylaod failed - %v\n", err)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, "http://192.168.1.10:60001/webhook", bytes.NewReader(json_data))
		if err != nil {
			fmt.Printf(":: Failed to initialize request - %v\n", err)
			continue
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf(":: Failed to connect to host - %v\n", err)
			continue
		}

		defer resp.Body.Close()

		body_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf(":: Failed to read response body - %v\n", err)
			continue
		}

		if string(body_bytes)[len(body_bytes)-1] == '\n' {
			body_bytes = body_bytes[:len(body_bytes)-1]
		}

		fmt.Printf("<<< [ %d ] %s\n", resp.StatusCode, string(body_bytes))

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf(":: Reading failrue - %v\n", err)
		return
	}
}

func parse_command(text string) (*Command, error) {
	if len(text) == 0 {
		return nil, errors.New("Empty command")
	}

	if text[0] != '/' {
		return nil, errors.New("Not a command. Commands starts with '/'")
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

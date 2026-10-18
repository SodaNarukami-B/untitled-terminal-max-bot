package server

import (
	"errors"
	"net/http"
	"strings"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func httpError(resw http.ResponseWriter, msg string, stat int) {
	resw.WriteHeader(stat)
	resw.Write([]byte(msg))
}

// Command structure needed to parse received text in updates that
// contains Message as command
type Command struct {
	Command string
	Args    []string
	RawArgs string
}

// ParseCommand parses text as Command structure
// "/comm arg1 arg2" -> &Command{Command: "comm", Args: ["arg1", "arg2"], RawArgs: "arg1 arg2",}
func ParseCommand(t string) (*Command, error) {
	if len(t) == 0 {
		return nil, errors.New("Empty string")
	}

	if t[0] != '/' {
		return nil, errors.New("Not a command")
	}

	t = t[1:]

	command, rawArgs, _ := strings.Cut(t, " ")

	args := strings.Fields(rawArgs)

	return &Command{
		Command: command,
		Args:    args,
		RawArgs: rawArgs,
	}, nil
}

type Attachment struct {
	T        string
	Url      string
	Filename string
	Size     int64
}

type MCU_ServerMessage struct {
	User        schemes.User `json:"user"`
	Timestamp   int64        `json:"timestamp"`
	Attachments []Attachment `json:"attachments"`
}

// MCU_Command_handler is a type of MessageCratedUpdate command hanlders. Used for handler tables.
type MCU_COMMAND_HANDLER func(http.ResponseWriter, *schemes.Message, Command, *http.Client)

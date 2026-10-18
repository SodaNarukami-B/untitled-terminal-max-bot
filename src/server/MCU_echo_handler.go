package server

import (
	"net/http"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func MCU_echo(resw http.ResponseWriter, m *schemes.Message, command Command, _ *http.Client) {
	resw.WriteHeader(http.StatusOK)
	resw.Write([]byte(command.RawArgs))
}

package server

import (
	"net/http"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func MCU_echo_handler(resw http.ResponseWriter, req *http.Request, msg *schemes.Message, command Command) {
	resw.WriteHeader(http.StatusOK)
	resw.Write([]byte(command.Raw_args))
}

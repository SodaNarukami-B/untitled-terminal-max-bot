package main

import (
	"net/http"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func handle_echo(resw http.ResponseWriter, req *http.Request, msg *schemes.Message) {
	resw.WriteHeader(http.StatusOK)

	comm, _ := parse_command(msg.Body.Text)
	resw.Write([]byte(comm.raw_args))
}

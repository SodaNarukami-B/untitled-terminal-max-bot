package main

import (
	"net/http"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func handle_echo(resw http.ResponseWriter, req *http.Request, msg *schemes.Message) {
	comm, _ := parse_command(msg.Body.Text)

	resw.WriteHeader(http.StatusOK)
	resw.Write([]byte(comm.raw_args))
}

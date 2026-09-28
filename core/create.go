// draft

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

type CreateMeta struct {
	UserId int64 `json:"user_id"`
	// first name and sec name ???
	Timestamp int64 `json:"timestamp"`
}

type CreateMessage struct {
	Meta        CreateMeta        `json:"meta"`
	Attachments []json.RawMessage `json:"attachments"`
}

func hanlde_create(resw http.ResponseWriter, req *http.Request, msg *schemes.Message) {
	var create_message CreateMessage

	create_message.Meta.UserId = msg.Sender.UserId
	create_message.Meta.Timestamp = time.Now().Unix() // XXX: idk, we maybe may use timestamp from msg
	create_message.Attachments = msg.Body.RawAttachments

	if len(create_message.Attachments) == 0 {
		resw.WriteHeader(http.StatusBadRequest)
		resw.Write([]byte("Invalid request"))

		return
	}

	payload, _ := json.Marshal(create_message) // XXX: i think you must to check err

	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	newreq, err := http.NewRequest("POST", "http://database.db/branches/master", bytes.NewReader(payload))

	if err != nil {
		resw.WriteHeader(http.StatusBadRequest)
		resw.Write([]byte("Invalid packet structure"))

		return
	}

	resp, err := client.Do(newreq)

	if err != nil {
		log.Fatalf("Failed to connect to database")

		resw.WriteHeader(http.StatusInternalServerError)
		resw.Write([]byte("Internal server error"))

		return
	}

	switch resp.StatusCode {
	case 201:
		{
			resw.WriteHeader(http.StatusCreated)
			resw.Write([]byte("File updloaded"))

			return
		}
	case 400:
		{
			resw.WriteHeader(http.StatusBadRequest)
			resw.Write([]byte("Failed to upload file"))

			return
		}

	default:
		{
			resw.WriteHeader(http.StatusInternalServerError)
			resw.Write([]byte("Unexpected exception"))

			return
		}
	}
}

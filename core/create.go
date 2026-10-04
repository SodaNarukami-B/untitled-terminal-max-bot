// draft

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Attachment struct {
	T        string `json:"type"`
	Url      string `json:"url"`
	Filename string `json:"filename"`
}

type CreateMessage struct {
	User        schemes.User `json:"user"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   int64        `json:"timestamp"`
}

var TYPE_HANDLERS = map[string]func(string){
	"file": func(t string) {},
}

func handle_create(resw http.ResponseWriter, req *http.Request, msg *schemes.Message) {

	var create_message CreateMessage

	// user
	create_message.User = msg.Sender

	// attachments
	for _, rawItem := range msg.Body.RawAttachments {
		var Item map[string]any

		err := json.Unmarshal(rawItem, &Item)

		if err != nil {
			http.Error(resw, "Invalid attachment", http.StatusBadRequest)
			return
		}

		attachment_type, exists := Item["type"].(string)
		if !exists {
			http.Error(resw, "Invalid attachment", http.StatusBadRequest)
			return
		}

		payload, exists := Item["payload"].(map[string]any)
		if !exists {
			http.Error(resw, "Invalid attachment", http.StatusBadRequest)
			return
		}

		// URL
		// 1
		host, exists := payload["host"].(string)
		if !exists {
			http.Error(resw, "Invalid attachment", http.StatusBadRequest)
			return
		}

		// 2
		token, exists := payload["token"].(string)
		if !exists {
			http.Error(resw, "Invalid attachment", http.StatusBadRequest)
			return
		}

		// 1 + 2

		url, err := url.JoinPath(host, token)
		if err != nil {
			http.Error(resw, "Invalid source", http.StatusBadRequest)
			return
		}

		switch schemes.AttachmentType(attachment_type) {
		case schemes.AttachmentFile:
			{
				filename, exists := Item["filename"].(string)
				if !exists {
					http.Error(resw, "Invalid attachment", http.StatusBadRequest)
					return
				}

				create_message.Attachments = append(create_message.Attachments, Attachment{
					T:        attachment_type,
					Url:      url,
					Filename: filename,
				})
			}
			// FIXME: Here MUST be all supported attachment types
		default:
			{
				http.Error(resw, "Invalid/not implemented attachment", http.StatusBadRequest)
				return
			}
		}
	}

	create_message.Timestamp = time.Now().Unix()

	if len(create_message.Attachments) == 0 {
		http.Error(resw, "Please, attach file to upload", http.StatusBadRequest)

		return
	}

	payload, err := json.Marshal(create_message)

	if err != nil {
		http.Error(resw, "Internal server error", http.StatusInternalServerError)
		return
	}

	//XXX: ------------- TESTING BLOCK -----------------------
	resw.WriteHeader(http.StatusOK)
	resw.Write([]byte("Accepted\n"))

	pretty_payload, _ := json.MarshalIndent(create_message, "json >>", "  ")

	log.Printf("Apporacing to send: \njson >>%s\n", string(pretty_payload))

	return
	//XXX: ------- ^^^ REMOVE AFTER TESTING ^^^-----------------

	// Database fetch

	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	newreq, err := http.NewRequest("POST", "http://database.db/branches/master", bytes.NewReader(payload))

	if err != nil {
		http.Error(resw, "Invalid packet structure", http.StatusBadRequest)
		return
	}

	resp, err := client.Do(newreq)

	if err != nil {
		log.Fatalf("Failed to connect to database")

		http.Error(resw, "Internal server error", http.StatusInternalServerError)
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
			http.Error(resw, "Failed to upload file", http.StatusBadRequest)
			return
		}

	default:
		{
			http.Error(resw, "Unexpected exception", http.StatusInternalServerError)
			return
		}
	}
}

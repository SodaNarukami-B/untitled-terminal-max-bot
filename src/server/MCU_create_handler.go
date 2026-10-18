package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func MCU_create(resw http.ResponseWriter, msg *schemes.Message, com Command, dbClient *http.Client) {

	var serverMessage MCU_ServerMessage

	serverMessage.User = msg.Sender

	for _, rawItem := range msg.Body.RawAttachments {
		var item map[string]any
		err := json.Unmarshal(rawItem, &item)
		if err != nil {
			httpError(resw, "Bad Attachment (marshal)", http.StatusBadRequest)
			return
		}

		var at_type string
		var payload map[string]any
		var payload_host string
		var payload_token string
		var at_url string
		var at_filename string
		var at_size int64

		var ex bool

		if at_type, ex = item["type"].(string); !ex {
			httpError(resw, "Bad Attachment (type)", http.StatusBadRequest)
			return
		}

		if payload, ex = item["payload"].(map[string]any); !ex {
			httpError(resw, "Bad Attachment (payload)", http.StatusBadRequest)
			return
		}

		if payload_host, ex = payload["host"].(string); !ex {
			httpError(resw, "Bad Attachment (host)", http.StatusBadRequest)
			return
		}

		if payload_token, ex = payload["token"].(string); !ex {
			httpError(resw, "Bad Attachment (token)", http.StatusBadRequest)
			return
		}

		switch schemes.AttachmentType(at_type) {
		case schemes.AttachmentFile:
			if at_filename, ex = item["filename"].(string); !ex {
				httpError(resw, "Bad Attachment (filename)", http.StatusBadRequest)
				return
			}

			if at_sizeRaw, ex := item["size"].(float64); !ex {
				httpError(resw, "Bad Attachment (size)", http.StatusBadRequest)
				return
			} else {
				at_size = int64(at_sizeRaw)
			}

		default:
			httpError(resw, "Bad Attachment (file type)", http.StatusBadRequest)
			return
		}

		at_url, err = url.JoinPath(payload_host, payload_token)
		if err != nil {
			httpError(resw, "Bad Attachment (url compare)", http.StatusBadRequest)
			return
		}

		serverMessage.Attachments = append(serverMessage.Attachments, Attachment{
			T:        at_type,
			Url:      at_url,
			Filename: at_filename,
			Size:     at_size,
		})
	}

	serverMessage.Timestamp = time.Now().Unix()

	if len(serverMessage.Attachments) == 0 {
		httpError(resw, "No Attachments", http.StatusBadRequest)
		return
	}

	reqPayload, err := json.Marshal(serverMessage)

	if err != nil {
		httpError(resw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	dbReq, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:60002/repo", bytes.NewReader(reqPayload))

	dbReq.Header.Add("Content-Type", "application/json")
	dbReq.Header.Add("Content-Length", string(len(reqPayload)))

	if err != nil {
		httpError(resw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res, err := dbClient.Do(dbReq)

	// XXX: make repeating with exponentaly increasing delay
	if err != nil {
		log.Printf("Failed to connect to database (%s) [%v]\n", "http://127.0.0.1:60002/repo", err.Error())

		httpError(resw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch res.StatusCode {
	case 201:
		resw.WriteHeader(http.StatusCreated)
		resw.Write([]byte("Created"))

		log.Printf("CREATE > database\t[ CREATED ]\n")

		return
	case 400:
		rawRes, err := io.ReadAll(res.Body)
		if err != nil {
			httpError(resw, "Bad Gateway", http.StatusBadGateway)
			log.Printf("CREATE > database\t[ 502 ]\n")
			return
		}

		var jsonRes map[string]string
		err = json.Unmarshal(rawRes, &jsonRes)

		if err != nil {
			httpError(resw, "Bad Gateway", http.StatusBadGateway)
			log.Printf("CREATE > database\t[ 502 ]\n")
			return
		}

		var details string
		var ex bool

		if details, ex = jsonRes["details"]; !ex {
			httpError(resw, "Bad Gateway", http.StatusBadGateway)
			log.Printf("CREATE > database\t[ 502 ]\n")
			return
		}

		resw.WriteHeader(http.StatusBadRequest)
		resw.Write([]byte(details))

		log.Printf("CREATE > database\t[ 400 ]\n")

		return

	case 500:
		httpError(resw, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("CREATE > database\t[ 500 ]\n")

		return
	default:
		httpError(resw, "Bad Gateway", http.StatusBadGateway)
		log.Printf("CREATE > database\t[ 502 ]\n")
		return
	}
}

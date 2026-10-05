package client

import (
	"encoding/json"

	schemes "github.com/max-messenger/max-bot-api-client-go/schemes"
)

func COM_create_message(msg *schemes.Message, raw_args string) error {
	data := []byte(raw_args)

	err := json.Unmarshal(data, &msg.Body.RawAttachments)

	if err != nil {
		return err
	}

	msg.Body.Text = "/create"

	return nil
}

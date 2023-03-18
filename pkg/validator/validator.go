package validator

import (
	"CallFrescoBot/pkg/messages"
	"errors"
)

func Validate(cmd string) (string, error) {
	if cmd == "" {
		return messages.UnsupportedMessageType, errors.New("unsupported message type")
	}
	if len([]rune(cmd)) < 4 {
		return messages.MessageIsTooShort, errors.New("message is too short")
	}

	return cmd, nil
}

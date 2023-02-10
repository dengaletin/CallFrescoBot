package commands

import "CallFrescoBot/pkg/messages"

type StatusCommand struct {
	Message string
}

func (cmd StatusCommand) RunCommand() string {
	return messages.StatusMsg
}

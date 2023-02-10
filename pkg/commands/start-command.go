package commands

import "CallFrescoBot/pkg/messages"

type StartCommand struct {
	Message string
}

func (cmd StartCommand) RunCommand() string {
	return messages.StartMsg
}

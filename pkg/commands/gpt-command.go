package commands

import (
	gpt "CallFrescoBot/Gpt"
)

type GptCommand struct {
	Message string
}

func (cmd GptCommand) RunCommand() string {
	return gpt.GetResponse(cmd.Message)
}

package commands

import (
	"CallFrescoBot/pkg/consts"
)

type StartCommand struct {
	Message string
}

func (cmd StartCommand) RunCommand() string {
	return consts.StartMsg
}

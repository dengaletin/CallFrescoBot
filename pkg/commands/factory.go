package commands

type Factory interface {
	RunCommand() ICommand
}

func GetCommand(cmd string) ICommand {
	switch cmd {
	default:
		return GptCommand{Message: cmd}
	case Start:
		return StartCommand{}
	case Status:
		return StatusCommand{}
	}
}

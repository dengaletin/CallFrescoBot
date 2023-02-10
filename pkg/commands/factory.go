package commands

type Factory interface {
	RunCommand() ICommand
}

func GetCommand(cmd string) (ICommand, error) {
	switch cmd {
	default:
		return GptCommand{Message: cmd}, nil
	case Start:
		return StartCommand{}, nil
	case Status:
		return StatusCommand{}, nil
	}
}

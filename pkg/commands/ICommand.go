package commands

type ICommand interface {
	RunCommand() string
	Common() string
}

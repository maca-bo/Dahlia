package command

import "fmt"

type BadStatusCode struct {
	statusCode int
}

func (e *BadStatusCode) Error() string {
	return fmt.Sprintf("BadStatusCode: %d", e.statusCode)
}

type UnimplCommandError struct {
	cmd Op
}

func (e *UnimplCommandError) Error() string {
	return fmt.Sprintf("Command %s is Unimplemented", e.cmd)
}

package command

import "fmt"

type BadStatusCode struct {
	statusCode int
}

func (e *BadStatusCode) Error() string {
	return fmt.Sprintf("BadStatusCode: %d", e.statusCode)
}

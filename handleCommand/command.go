package command

import "fmt"

type Op int

const (
	Unlock Op = 0
	Lock   Op = 1
)

func (o Op) String() string {
	switch o {
	case Unlock:
		return "Unlock"
	case Lock:
		return "Lock"
	default:
		return fmt.Sprintf("%d", o)
	}
}

type Command struct {
	Operate  Op `json:"operate"`
	UserID   string
	Argument interface{} `json:"argument"`
}

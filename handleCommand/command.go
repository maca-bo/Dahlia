package command

type Op int

const (
	Unlock Op = 0
	Lock   Op = 1
)

type Command struct {
	Operate  Op `json:"operate"`
	UserID   string
	Argument interface{} `json:"argument"`
}

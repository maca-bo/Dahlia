package verify

import (
	"fmt"
	"time"
)

type OldTimeStumpError struct {
	t time.Duration
}

func (e *OldTimeStumpError) Error() string {
	return fmt.Sprintf("OldTimeError: %.0fm", e.t.Minutes())
}

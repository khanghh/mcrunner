package mcagent

import "fmt"

var (
	ErrTicketNotFound  = fmt.Errorf("ticket not found")
	ErrTicketExpired   = fmt.Errorf("ticket expired")
	ErrServiceMismatch = fmt.Errorf("service mismatch")
)

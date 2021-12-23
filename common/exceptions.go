package common

import (
	"fmt"
)

// NotFoundError an error indicates something is not found in the ledger
type NotFoundError struct {
	What string
}

// Error get the error message
func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.What)
}

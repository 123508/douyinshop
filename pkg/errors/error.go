package errors

import "fmt"

type BasicMessageError struct {
	Message string
}

func (e *BasicMessageError) Error() string {
	return fmt.Sprintf("error message: %s", e.Message)
}

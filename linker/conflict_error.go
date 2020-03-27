package linker

import "fmt"

// ConflictError represents a group of errors considered as a conflict.
type ConflictError struct {
	Errs []error
}

func (e *ConflictError) Error() string {
	errlen := len(e.Errs)
	verb, conflicts := resolveWords(errlen)
	return fmt.Sprintf("linker: there %s %d %s", verb, errlen, conflicts)
}

func resolveWords(n int) (string, string) {
	switch {
	case n == 1:
		return "is", "conflict"
	default:
		return "are", "conflicts"
	}
}

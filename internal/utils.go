package internal

// NewBool returns the address of b.
func NewBool(b bool) *bool { return &b }

// NewString returns the address of s.
func NewString(s string) *string { return &s }

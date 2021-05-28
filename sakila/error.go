package sakila

// Error is a sakila error.
type Error string

const (
	// ErrorInternal is an internal error.
	ErrorInternal = Error("internal")
	// ErrorNotFound is a resource not found error.
	ErrorNotFound = Error("not found")
)

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

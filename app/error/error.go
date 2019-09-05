package error

// Error is a custom error
type Error interface {
	error
	Status() int
}

// StatusError contains HTTP status error
// StatusError implements type Error interface
type StatusError struct {
	Err  error
	Code int
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status return status error code
func (se StatusError) Status() int {
	return se.Code
}

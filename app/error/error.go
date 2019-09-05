package error

// Error is a custom error
// type Error interface {
// 	error
// 	Status() int
// }

// StatusError contains HTTP status error
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status return status error code
func (se StatusError) Status() int {
	return se.Code
}

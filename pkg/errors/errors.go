package errors

type InternalError struct {
	Details string
}

func (err InternalError) Error() string {
	return err.Details
}

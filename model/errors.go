package model

type GenericError struct {
	Err error
}

func (e *GenericError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Err     error
	Message string
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message + ". " + e.Err.Error()
	}
	return e.Err.Error()
}

type ResourceDoesNotExistError struct {
	Err error
}

func (e *ResourceDoesNotExistError) Error() string {
	return e.Err.Error()
}

type AuthenticationError struct {
	Err error
}

func (e *AuthenticationError) Error() string {
	return e.Err.Error()
}

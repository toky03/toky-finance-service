package model

type BusinessError struct {
	errType string
	Err     error
	Cause   string
}

type TechnicalError struct {
	errType    string
	Err        error
	Stacktrace string
}

type TokyError interface {
	Error() error
	ErrorMessage() string
	IsTechnicalError() bool
	ErrorType() string
}

func CreateBusinessError(cause string, err error) BusinessError {
	return BusinessError{
		errType: "BusinessError",
		Err:     err,
		Cause:   cause,
	}
}

func CreateBusinessErrorNotFound(cause string, err error) BusinessError {
	return BusinessError{
		errType: "NotFoundError",
		Err:     err,
		Cause:   cause,
	}
}

func CreateBusinessValidationError(cause string, err error) BusinessError {
	return BusinessError{
		errType: "ValidationError",
		Err:     err,
		Cause:   cause,
	}
}

func CreateTechnicalError(stacktrace string, err error) TechnicalError {
	return TechnicalError{
		errType:    "TechnicalError",
		Err:        err,
		Stacktrace: stacktrace,
	}
}
func IsExisting(err TokyError) bool {
	return err != nil && err.Error() != nil
}

func IsExistingNotFoundError(err TokyError) bool {
	return err != nil && err.ErrorType() == "NotFoundError"
}

func IsExistingValidationError(err TokyError) bool {
	return err != nil && err.ErrorType() == "ValidationError"
}
func IsExistingBuisnessError(err TokyError) bool {
	return err != nil && err.ErrorType() == "BusinessError"
}

func (e TechnicalError) IsTechnicalError() bool {
	return true
}
func (e TechnicalError) Error() error {
	return e.Err
}
func (e TechnicalError) ErrorMessage() string {
	return e.Stacktrace
}

func (e TechnicalError) ErrorType() string {
	return e.errType
}

func (e BusinessError) IsTechnicalError() bool {
	return false
}
func (e BusinessError) Error() error {
	return e.Err
}

func (e BusinessError) ErrorType() string {
	return e.errType
}
func (e BusinessError) ErrorMessage() string {
	return e.Cause
}

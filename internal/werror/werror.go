package werror

type WError struct {
	Code    *int
	Message *string
	Err     error
}

func New() *WError {
	return NewFromError(nil)
}

func NewFromError(err error) *WError {
	return &WError{
		Code:    nil,
		Message: nil,
		Err:     err,
	}
}

func (err *WError) Error() string {
	if err.Message != nil {
		return *err.Message
	}

	return err.Error()
}

func (err *WError) SetMessage(msg string) *WError {
	err.Message = &msg
	return err
}

func (err *WError) SetCode(code int) *WError {
	err.Code = &code
	return err
}

func (err *WError) GetCodeOr(code int) int {
	if err.Code != nil {
		return *err.Code
	}

	return code
}

func (err *WError) GetMessageOr(msg string) string {
	if err.Message != nil {
		return *err.Message
	}

	return msg
}

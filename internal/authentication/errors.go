package authentication

type EmailFormatError struct {
	Email string
}

func (e EmailFormatError) Error() string {
	return e.Email
}

type PasswordMismatchError struct {
	Message string
}

func (e PasswordMismatchError) Error() string {
	return e.Message
}

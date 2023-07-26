package apperrors

import "errors"

var (
	ErrIncorrectVerificationCode = errors.New("incorrect verificaion code")
	ErrUserAlreadyExists         = errors.New("user with such email lready exists")
	ErrIncorrectUserData         = errors.New("first name, last name, company name ust be more than 2 characters long")
	ErrUserAlreadyVerifyed       = errors.New("user with such emailaready verifyed")
	ErrUserNotFound              = errors.New("user doesn't exists")
	ErrVerificationCodeExpired   = errors.New("The code has expired")
	ErrInternalServerError       = errors.New("Thre was an internal serr bug, please try again later")
	ErrIncorrectEmailFormat      = errors.New("Incorrect email format")
	ErrIncorrectPasswordFormat   = errors.New("Password must be at least 8 characters long and contain at least one uppercase letter and one digit")
	ErrUserBlocked               = errors.New("user is blocked")
	ErrDocumentNotFound          = errors.New("document doesn't exists")
)

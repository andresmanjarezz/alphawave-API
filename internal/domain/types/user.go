package types

import "time"

type UserSignUpDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	JobTitle  string `json:"jobTitle"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserSignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerificationCodeDTO struct {
	Email                       string        `json:"email"`
	VerificationCodeExpiresTime time.Duration `json:"verificationCodeExpiresTime"`
}

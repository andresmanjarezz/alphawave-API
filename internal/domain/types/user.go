package types

import "time"

type UserDTO struct {
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	JobTitle       string    `json:"jobTitle"`
	Email          string    `json:"email"`
	Verification   bool      `json:"verification" bson:"verification"`
	RegisteredTime time.Time `json:"registeredTime" bson:"registeredTime"`
	LastVisitTime  time.Time `json:"lastVisitTime" bson:"lastVisitTime"`
	Blocked        bool      `json:"blocked" bson:"blocked"`
}

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

type ForgotPasswordPayloadDTO struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	ResultToken string `json:"resultToken"`
}

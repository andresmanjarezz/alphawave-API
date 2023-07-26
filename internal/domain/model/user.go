package model

import "time"

type User struct {
	ID             string                  `json:"id" bson:"_id,omitempty"`
	FirstName      string                  `json:"firstName" bson:"firstName"`
	LastName       string                  `json:"lastName" bson:"lastName"`
	CompanyName    string                  `json:"companyName" bson:"companyName"`
	Email          string                  `json:"email" bson:"email"`
	Password       string                  `json:"password" bson:"password"`
	Verification   UserVerificationPayload `json:"verification" bson:"verification"`
	RegisteredTime time.Time               `json:"registeredTime" bson:"registeredTime"`
	LastVisitTime  time.Time               `json:"lastVisitTime" bson:"lastVisitTime"`
	Blocked        bool                    `json:"blocked" bson:"blocked"`
}

type UserVerificationPayload struct {
	VerificationCode            string    `json:"varificationCode" bson:"verificationCode"`
	VerificationCodeExpiresTime time.Time `json:"verificationCodeExpiresTime" bson:"verificationCodeExpiresTime"`
	Verified                    bool      `json:"verified" bson:"verified"`
}

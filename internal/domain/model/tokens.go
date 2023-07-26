package model

import "time"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresTime  time.Time `json:"expiresTime" bson:"expiresTime"`
}

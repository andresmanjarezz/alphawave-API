package model

import "time"

const (
	OrderStatusCreated = "created"
	OrderStatusPaid    = "paid"
	OrderStatusFailed  = "failed"
)

type Subscription struct {
	ID           string        `json:"id" bson:"_id,omitempty"`
	TeamID       string        `json:"teamID" bson:"teamID"`
	UserInfo     UserInfoShort `json:"userInfo" bson:"userInfo"`
	Amount       string        `json:"amount" bson:"amount"`
	Currency     string        `json:"currency" bson:"currency"`
	ExpiresTime  time.Time     `json:"expiresTime" bson:"expiresTime"`
	Transactions []Transaction `json:"transaction" bson:"transaction"`
	Status       string        `json:"status" bson:"status"`
}

type Transaction struct {
	Status       string    `json:"status" bson:"status"`
	CreatedDate  time.Time `json:"createdDate" bson:"createdDate"`
	Description  string    `json:"description" bson:"description"`
	InvoiceTotal string    `json:"invoiceTotal" bson:"invoiceTotal"`
}

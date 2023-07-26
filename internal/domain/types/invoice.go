package types

import "time"

type InvoiceDTO struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userID"`
	InvTitle    string    `json:"invTitle"`
	InvNum      int       `json:"invNum"`
	CreatedDate time.Time `json:"createdDate"`
	Balance     float64   `json:"balance"`
	Notes       string    `json:"notes"`
	Dispatch    bool      `json:"dispatch"`
	Discount    bool      `json:"discount"`
	ColorLine   string    `json:"colorLine"`
	Currency    string    `json:"currency"`
	From        From      `json:"from"`
	To          To        `json:"to"`
	IsMarked    bool      `json:"isMarked"`
	InvList     []InvItem `json:"invList"`
}

type From struct {
	Name          string  `json:"name"`
	EmailFrom     string  `json:"emailFrom"`
	Address       Address `json:"address"`
	Phone         string  `json:"phone"`
	BusinessPhone string  `json:"businessPhone"`
}

type To struct {
	Name    string  `json:"name"`
	EmailTo string  `json:"emailTo"`
	Address Address `json:"address"`
	Phone   string  `json:"phone"`
}

type Address struct {
	Street    string `json:"street"`
	CityState string `json:"cityState"`
	ZipCode   string `json:"zipCode"`
}

type InvItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Qty         float64 `json:"qty"`
	Amount      float64 `json:"amount"`
}

type InvItemAmountDTO struct {
	Amount float64 `json:"amount"`
}

package model

type Package struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Price       uint   `json:"price" bson:"price"`
	Currency    string `json:"currency" bson:"currency"`
}

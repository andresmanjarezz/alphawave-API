package model

type Team struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	TeamName string `json:"teamName" bson:"teamName"`
	JobTitle string `json:"jobTitle" bson:"jobTitle"`
	OwnerID  string `json:"ownerID" bson:"ownerID"`
}

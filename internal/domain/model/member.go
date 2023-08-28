package model

type Member struct {
	ID     string   `json:"id" bson:"_id,omitempty"`
	TeamID string   `json:"teamID" bson:"teamID"`
	UserID string   `json:"userID" bson:"userID"`
	Status string   `json:"status" bson:"status"`
	Roles  []string `json:"roles" bson:"roles"`
}

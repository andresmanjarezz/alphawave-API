package model

type TeamRoles struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	TeamID string `json:"teamID" bson:"teamID"`
	Roles  []Role `json:"roles" bson:"roles"`
}

type Role struct {
	Role        string   `json:"role" bson:"role"`
	Permissions []string `json:"permissions" bson:"permissions"`
}

package model

type Task struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	UserID   string `json:"userID" bson:"userID"`
	Title    string `json:"title" bson:"title"`
	Status   string `json:"status" bson:"status"`
	Priority string `json:"priority" bson:"priority"`
	Order    int    `json:"order" bson:"order"`
}

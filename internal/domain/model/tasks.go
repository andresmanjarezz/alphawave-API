package model

type Task struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title"`
}

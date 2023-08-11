package types

type TaskDTO struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title"`
	Order int    `json:"order" bson:"order"`
}

type TasksCreateDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

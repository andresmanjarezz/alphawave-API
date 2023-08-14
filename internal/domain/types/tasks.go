package types

type TaskDTO struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Order    int    `json:"order"`
}

type TasksCreateDTO struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Order    int    `json:"order"`
}

type UpdateTaskDTO struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Order    int    `json:"order"`
}

package models

type Job struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	Skills      []string `json:"skills"`
	Salary      string   `json:"salary"`
	Description string   `json:"description"`
}

type BindJob struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	Skills      []string `json:"skills"`
	Salary      string   `json:"salary"`
	Description string   `json:"description"`
}

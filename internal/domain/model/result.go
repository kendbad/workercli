package model

type Result struct {
	TaskID  string
	Status  string
	Details string
	Error   string
}

type TaskResult struct {
	TaskID string
	Status string
}

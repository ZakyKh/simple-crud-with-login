package tasks

import(
)

type TaskRequest struct {
	Description string `json:"description"`
	Difficulty  int    `json:"difficulty"`
	Done        bool   `json:"done"`
}

type Task struct {
	Id          int    `                   json:"id"`
	Description string `gorm:"description" json:"description"`
	Difficulty  int    `gorm:"difficulty"  json:"difficulty"`
	Done        bool   `gorm:"done"        json:"done"`
}
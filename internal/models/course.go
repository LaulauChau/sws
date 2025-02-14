package models

type Course struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Date  string `json:"date"`
	Start string `json:"start"`
	End   string `json:"end"`
}

package models

type Book struct {
	ID      uint   `json:"id,omitempty" gorm:"primary_key"`
	Title   string `json:"title,omitempty" gorm:"uniqueIndex"`
	Author  string `json:"author,omitempty"`
	Summary string `json:"summary,omitempty"`
}

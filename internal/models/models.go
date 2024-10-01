package models

type Song struct {
	ID          int    `json:"id" db:"id"`
	Group       string `json:"group" db:"group"`
	Song        string `json:"song" db:"song"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}

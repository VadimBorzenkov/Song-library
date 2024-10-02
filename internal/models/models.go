package models

type Song struct {
	ID          int    `json:"id" db:"id"`
	Group       string `json:"group" db:"group_name"`
	Song        string `json:"song" db:"song_name"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}

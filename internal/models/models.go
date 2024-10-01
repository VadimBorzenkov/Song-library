package models

// Song представляет песню в базе данных
type Song struct {
	ID          int    `json:"id" db:"id"`                     // Уникальный идентификатор
	Group       string `json:"group" db:"group"`               // Название группы
	Song        string `json:"song" db:"song"`                 // Название песни
	ReleaseDate string `json:"release_date" db:"release_date"` // Дата выпуска
	Text        string `json:"text" db:"text"`                 // Текст песни
	Link        string `json:"link" db:"link"`                 // Ссылка на песню
}

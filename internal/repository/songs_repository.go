package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/sirupsen/logrus"
)

func (repo *ApiRepository) GetData(filter map[string]string, limit int, offset int) ([]models.Song, error) {
	var songs []models.Song
	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE 1=1"
	args := []interface{}{}

	// Добавляем фильтрацию по полям
	i := 1
	for key, value := range filter {
		query += fmt.Sprintf(" AND %s = $%d", key, i)
		args = append(args, value)
		i++
	}

	// Добавляем параметры LIMIT и OFFSET
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	// Логируем SQL-запрос
	repo.logger.WithFields(logrus.Fields{
		"query": query,
		"args":  args,
	}).Debug("Executing GetData query")

	// Выполняем запрос
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		repo.logger.Error("Error executing GetData query: ", err)
		return nil, err
	}
	defer rows.Close()

	// Сканируем результаты
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			repo.logger.Error("Error scanning GetData rows: ", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	repo.logger.Infof("Successfully fetched %d songs", len(songs))
	return songs, nil
}

func (repo *ApiRepository) GetSongPagi(id int, limit int, offset int) (*models.Song, error) {
	var song models.Song
	err := repo.db.QueryRow("SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = $1", id).Scan(
		&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link,
	)
	if err != nil {
		repo.logger.Error("Error fetching song for pagination: ", err)
		return nil, err
	}

	// Логируем извлеченные данные о песне
	repo.logger.WithFields(logrus.Fields{
		"songID":   id,
		"songName": song.Song,
		"group":    song.Group,
	}).Debug("Successfully fetched song")

	// Разбиваем текст на куплеты по новой строке и применяем пагинацию
	verses := strings.Split(song.Text, "\n")
	if offset >= len(verses) {
		repo.logger.Warn("Offset out of range for GetSongPagi")
		return nil, errors.New("offset out of range")
	}

	end := offset + limit
	if end > len(verses) {
		end = len(verses)
	}

	song.Text = strings.Join(verses[offset:end], "\n")
	repo.logger.Infof("Returning %d verses from song '%s'", end-offset, song.Song)
	return &song, nil
}

func (r *ApiRepository) DeleteSong(id int) error {
	_, err := r.db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		r.logger.Error("Error deleting song: ", err)
		return err
	}

	r.logger.Infof("Song with ID %d successfully deleted", id)
	return nil
}

func (r *ApiRepository) UpdateSongData(song *models.Song) error {
	_, err := r.db.Exec(`UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`,
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID,
	)
	if err != nil {
		r.logger.Error("Error updating song: ", err)
		return err
	}

	r.logger.Infof("Song with ID %d successfully updated", song.ID)
	return nil
}

func (r *ApiRepository) AddNewSong(song *models.Song) error {
	_, err := r.db.Exec(`INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`,
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link,
	)
	if err != nil {
		r.logger.Error("Error inserting new song: ", err)
		return err
	}

	r.logger.Infof("New song '%s' by group '%s' added successfully", song.Song, song.Group)
	return nil
}

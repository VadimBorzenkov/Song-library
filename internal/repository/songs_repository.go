package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/VadimBorzenkov/online-song-library/internal/models"
)

func (repo *ApiRepository) GetData(filter map[string]string, limit int, offset int) ([]models.Song, error) {
	var songs []models.Song
	query := "SELECT id, group_name, song_name, COALESCE(release_date::text, '') AS release_date, COALESCE(text, '') AS text, COALESCE(link, '') AS link FROM songs WHERE 1=1"
	args := []interface{}{}

	i := 1
	for key, value := range filter {
		if key == "release_date" {
			query += fmt.Sprintf(" AND release_date = $%d", i)
		} else {
			query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
		}
		args = append(args, value)
		i++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		repo.logger.Error("Error executing GetData query: ", err)
		return nil, err
	}
	defer rows.Close()

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

	verses := strings.Split(song.Text, "\n\n")

	if offset >= len(verses) {
		repo.logger.Warn("Offset out of range for GetSongPagi")
		return nil, errors.New("offset out of range")
	}

	end := offset + limit
	if end > len(verses) {
		end = len(verses)
	}

	song.Text = strings.Join(verses[offset:end], "\n\n")
	repo.logger.Infof("Returning %d verses from song '%s'", end-offset, song.Song)
	return &song, nil
}

func (r *ApiRepository) DeleteSong(id int) (int64, error) {
	result, err := r.db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		r.logger.Error("Error deleting song: ", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Error fetching rows affected: ", err)
		return 0, err
	}

	if rowsAffected == 0 {
		r.logger.Warnf("No song found with ID %d", id)
	}

	return rowsAffected, nil
}

func (r *ApiRepository) UpdateSongData(song *models.Song) error {
	query := `UPDATE songs SET`
	params := []interface{}{}
	paramCounter := 1

	if song.Group != "" {
		query += ` group_name = $` + strconv.Itoa(paramCounter) + `,`
		params = append(params, song.Group)
		paramCounter++
	}

	if song.Song != "" {
		query += ` song_name = $` + strconv.Itoa(paramCounter) + `,`
		params = append(params, song.Song)
		paramCounter++
	}

	if song.Text != "" {
		query += ` text = $` + strconv.Itoa(paramCounter) + `,`
		params = append(params, song.Text)
		paramCounter++
	}

	if song.Link != "" {
		query += ` link = $` + strconv.Itoa(paramCounter) + `,`
		params = append(params, song.Link)
		paramCounter++
	}

	if song.ReleaseDate != "" {
		query += ` release_date = $` + strconv.Itoa(paramCounter) + `,`
		params = append(params, song.ReleaseDate)
		paramCounter++
	}

	if len(params) == 0 {
		r.logger.Error("No fields to update")
		return errors.New("no fields to update")
	}

	query = query[:len(query)-1]

	query += ` WHERE id = $` + strconv.Itoa(paramCounter)
	params = append(params, song.ID)

	_, err := r.db.Exec(query, params...)
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

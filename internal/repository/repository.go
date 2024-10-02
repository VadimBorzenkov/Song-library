package repository

import (
	"database/sql"

	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetData(filter map[string]string, limit int, offset int) ([]models.Song, error)
	GetSongPagi(id int, limit int, offset int) (*models.Song, error)
	DeleteSong(id int) (int64, error)
	UpdateSongData(song *models.Song) error
	AddNewSong(song *models.Song) error
}

type ApiRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewApiRepository(db *sql.DB, logger *logrus.Logger) *ApiRepository {
	return &ApiRepository{
		db:     db,
		logger: logger,
	}
}

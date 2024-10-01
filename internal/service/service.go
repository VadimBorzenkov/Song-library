package service

import (
	"github.com/VadimBorzenkov/online-song-library/config"
	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/VadimBorzenkov/online-song-library/internal/repository"
	externalapi "github.com/VadimBorzenkov/online-song-library/pkg/external_api"
	"github.com/sirupsen/logrus"
)

const (
	DateFormat = "2006-01-02 15:04:05"
)

type SongService interface {
	GetSongsWithPaginate(filter map[string]string, limit, offset int) ([]models.Song, error)
	GetSongWithVerses(id, limit, offset int) (*models.Song, error)
	AddNewSong(group, song string) (*models.Song, error)
	DeleteSong(id int) error
	UpdateSong(song *models.Song) error
}

type ApiService struct {
	repo   repository.Repository
	logger *logrus.Logger
	cfg    *config.Config
	exApi  *externalapi.ExternalApiClient
}

func NewApiService(repo repository.Repository, logger *logrus.Logger, cfg *config.Config) *ApiService {
	client := externalapi.NewExternalApiClient(cfg.ExternalApiURL, logger)
	return &ApiService{
		repo:   repo,
		logger: logger,
		cfg:    cfg,
		exApi:  client,
	}
}

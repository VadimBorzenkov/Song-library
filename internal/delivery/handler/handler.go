package handler

import (
	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/VadimBorzenkov/online-song-library/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	GetSongs(ctx *fiber.Ctx) error
	GetSongWithVerses(ctx *fiber.Ctx) error
	DeleteSong(ctx *fiber.Ctx) error
	AddNewSong(ctx *fiber.Ctx) error
	UpdateSong(ctx *fiber.Ctx) error
}

type CommonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type DataResponseSongs struct {
	Data    []models.Song `json:"data"`
	Message string        `json:"message"`
}

type DataResponseSong struct {
	Data    *models.Song `json:"data"`
	Message string       `json:"message"`
}

type ApiHandler struct {
	serv   service.SongService
	logger *logrus.Logger
}

func NewApiHandler(serv service.SongService, logger *logrus.Logger) *ApiHandler {
	return &ApiHandler{serv: serv, logger: logger}
}

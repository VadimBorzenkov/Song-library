package handler

import (
	"strconv"

	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type request struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

func (h *ApiHandler) GetSongs(ctx *fiber.Ctx) error {
	filters := make(map[string]string)
	for _, key := range []string{"group", "song", "releaseDate", "text", "link"} {
		value := ctx.Query(key)
		if value != "" {
			filters[key] = value
		}
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10")) // Значение по умолчанию 10
	if err != nil || limit <= 0 {
		h.logger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid limit value")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Limit must be a positive integer",
		})
	}

	page, err := strconv.Atoi(ctx.Query("page", "1")) // Значение по умолчанию 1
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid page value")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Invalid page number",
		})
	}

	offset := (page - 1) * limit

	h.logger.WithFields(logrus.Fields{
		"filters": filters,
		"limit":   limit,
		"page":    page,
		"offset":  offset,
	}).Debug("Fetching songs with filters")

	songs, err := h.serv.GetSongsWithPaginate(filters, limit, offset)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"filters": filters,
			"limit":   limit,
			"offset":  offset,
			"error":   err,
		}).Error("Error fetching songs")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Failed to fetch songs",
		})
	}

	h.logger.WithFields(logrus.Fields{
		"count": len(songs),
		"page":  page,
		"limit": limit,
	}).Info("Songs fetched successfully")

	return ctx.JSON(fiber.Map{
		"data":    songs,
		"message": "Songs retrieved successfully",
		"page":    page,
		"limit":   limit,
	})
}

func (h *ApiHandler) GetSongWithVerses(ctx *fiber.Ctx) error {
	songID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		h.logger.WithField("error", err).Warn("Invalid song ID")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Song ID must be a valid integer",
		})
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "5"))
	if err != nil || limit <= 0 {
		h.logger.WithField("error", err).Warn("Invalid limit for verses")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Limit must be a positive integer",
		})
	}

	offset, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil || offset < 0 {
		h.logger.WithField("error", err).Warn("Invalid offset value")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid offset value",
			"message": "Offset must be a non-negative integer",
		})
	}

	h.logger.WithFields(logrus.Fields{
		"songID": songID,
		"limit":  limit,
		"offset": offset,
	}).Debug("Fetching song with verses")

	song, err := h.serv.GetSongWithVerses(songID, limit, offset)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"songID": songID,
			"limit":  limit,
			"offset": offset,
			"error":  err,
		}).Error("Error fetching song with verses")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to fetch song with verses",
		})
	}

	h.logger.WithFields(logrus.Fields{
		"song": song.Song,
		"ID":   song.ID,
	}).Info("Song with verses fetched successfully")

	return ctx.JSON(fiber.Map{
		"data":    song,
		"message": "Song with verses retrieved successfully",
	})
}

func (h *ApiHandler) DeleteSong(ctx *fiber.Ctx) error {
	songID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		h.logger.WithField("error", err).Warn("Invalid song ID")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid song ID",
			"message": "Song ID must be a valid integer",
		})
	}

	h.logger.WithField("songID", songID).Info("Deleting song")

	if err := h.serv.DeleteSong(songID); err != nil {
		h.logger.WithField("songID", songID).Error("Error deleting song")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete song",
			"message": err.Error(),
		})
	}

	h.logger.WithField("songID", songID).Info("Song deleted successfully")
	return ctx.JSON(fiber.Map{
		"message": "Song deleted successfully",
	})
}

func (h *ApiHandler) UpdateSong(ctx *fiber.Ctx) error {
	songID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		h.logger.WithField("error", err).Warn("Invalid song ID")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid song ID",
			"message": "Song ID must be a valid integer",
		})
	}

	var songData models.Song
	if err := ctx.BodyParser(&songData); err != nil {
		h.logger.WithField("error", err).Warn("Failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": "Failed to parse JSON",
		})
	}

	songData.ID = songID

	h.logger.WithFields(logrus.Fields{
		"songID": songID,
		"song":   songData.Song,
		"group":  songData.Group,
	}).Info("Updating song")

	if err := h.serv.UpdateSong(&songData); err != nil {
		h.logger.WithField("songID", songID).Error("Error updating song")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update song",
			"message": err.Error(),
		})
	}

	h.logger.WithField("songID", songID).Info("Song updated successfully")
	return ctx.JSON(fiber.Map{
		"message": "Song updated successfully",
		"data":    songData,
	})
}

func (h *ApiHandler) AddNewSong(ctx *fiber.Ctx) error {
	var req request
	if err := ctx.BodyParser(&req); err != nil {
		h.logger.WithField("error", err).Warn("Failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	if req.Group == "" || req.Song == "" {
		h.logger.Warn("Group and Song fields are required")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Group and Song fields are required",
			"message": "Please provide both group and song names",
		})
	}

	h.logger.WithFields(logrus.Fields{
		"group": req.Group,
		"song":  req.Song,
	}).Info("Adding new song")

	newSong, err := h.serv.AddNewSong(req.Group, req.Song)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"group": req.Group,
			"song":  req.Song,
			"error": err,
		}).Error("Error adding new song")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to add new song",
			"message": err.Error(),
		})
	}

	h.logger.WithFields(logrus.Fields{
		"group": newSong.Group,
		"song":  newSong.Song,
	}).Info("New song added successfully")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    newSong,
		"message": "Song added successfully",
	})
}

package handler

import (
	"fmt"
	"strconv"

	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type request struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// GetSongs retrieves a list of songs based on optional filters, pagination, and limit.
// @Summary Get songs
// @Description Fetches a list of songs with optional filters and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group name"
// @Param song query string false "Filter by song name"
// @Param releaseDate query string false "Filter by release date"
// @Param text query string false "Filter by text content"
// @Param link query string false "Filter by link"
// @Param limit query int false "Number of results to return (default is 10)"
// @Param page query int false "Page number for pagination (default is 1)"
// @Success 200 {object} DataResponseSongs
// @Failure 400 {object} ErrorResponse
// @Router /songs/ [get]
func (h *ApiHandler) GetSongs(ctx *fiber.Ctx) error {
	filters := make(map[string]string)

	if group := ctx.Query("group"); group != "" {
		filters["group_name"] = group
	}
	if song := ctx.Query("song"); song != "" {
		filters["song_name"] = song
	}
	if releaseDate := ctx.Query("releaseDate"); releaseDate != "" {
		filters["release_date"] = releaseDate
	}

	for _, key := range []string{"release_date", "text", "link"} {
		value := ctx.Query(key)
		if value != "" {
			filters[key] = value
		}
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit <= 0 {
		h.logger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid limit value")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Limit must be a positive integer",
		})
	}

	page, err := strconv.Atoi(ctx.Query("page", "1"))
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
	if offset < 0 {
		h.logger.WithField("error", err).Warn("Invalid offset value")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid offset value",
			Message: "Offset must be a non-negative integer",
		})
	}

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

	return ctx.JSON(DataResponseSongs{
		Data:    songs,
		Message: "Songs retrieved successfully",
	})
}

// GetSongWithVerses retrieves a specific song by its ID along with its verses.
// @Summary Get song with verses
// @Description Fetches a specific song along with its verses by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param limit query int false "Number of verses to return (default is 5)"
// @Param offset query int false "Offset for verses (default is 0)"
// @Success 200 {object} DataResponseSong
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/get_song/{id} [get]
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
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid offset value",
			Message: "Offset must be a non-negative integer",
		})
	}

	song, err := h.serv.GetSongWithVerses(songID, limit, offset)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"songID": songID,
			"limit":  limit,
			"offset": offset,
			"error":  err,
		}).Error("Error fetching song with verses")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Failed to fetch song with verses",
		})
	}

	return ctx.JSON(DataResponseSong{
		Data:    song,
		Message: "Song with verses retrieved successfully",
	})
}

// DeleteSong removes a specific song by its ID.
// @Summary Delete song
// @Description Deletes a specific song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/delete_song/{id} [delete]
func (h *ApiHandler) DeleteSong(ctx *fiber.Ctx) error {
	songID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		h.logger.WithField("error", err).Warn("Invalid song ID")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid song ID",
			Message: "Song ID must be a valid integer",
		})
	}

	h.logger.WithField("songID", songID).Info("Deleting song")

	err = h.serv.DeleteSong(songID)
	if err != nil {
		if err.Error() == fmt.Sprintf("song with ID %d not found", songID) {
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Song not found",
				Message: fmt.Sprintf("Song with ID %d not found", songID),
			})
		}

		h.logger.WithField("songID", songID).Error("Error deleting song")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to delete song",
			Message: err.Error(),
		})
	}

	return ctx.JSON(SuccessResponse{
		Message: "Song deleted successfully",
	})
}

// UpdateSong modifies the details of a specific song by its ID.
// @Summary Update song
// @Description Updates a specific song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song data"
// @Success 200 {object} DataResponseSong
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/update_song/{id} [put]
func (h *ApiHandler) UpdateSong(ctx *fiber.Ctx) error {
	songID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		h.logger.WithField("error", err).Warn("Invalid song ID")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid song ID",
			Message: "Song ID must be a valid integer",
		})
	}

	var songData models.Song
	if err := ctx.BodyParser(&songData); err != nil {
		h.logger.WithField("error", err).Warn("Failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse JSON",
		})
	}

	songData.ID = songID

	if songData.Song == "" && songData.Group == "" && songData.Text == "" && songData.Link == "" && songData.ReleaseDate == "" {
		h.logger.WithField("songID", songID).Error("No fields to update")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to update song",
			Message: "no fields to update",
		})
	}

	if err := h.serv.UpdateSong(&songData); err != nil {
		h.logger.WithField("songID", songID).Error("Error updating song")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to update song",
			Message: err.Error(),
		})
	}

	return ctx.JSON(DataResponseSong{
		Data:    &songData,
		Message: "Song updated successfully",
	})
}

// AddNewSong creates a new song entry based on the provided request data.
// @Summary Add new song
// @Description Adds a new song to the library
// @Tags songs
// @Accept json
// @Produce json
// @Param request body request true "New song request"
// @Success 201 {object} DataResponseSong
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/add_song [post]
func (h *ApiHandler) AddNewSong(ctx *fiber.Ctx) error {
	var req request
	if err := ctx.BodyParser(&req); err != nil {
		h.logger.WithField("error", err).Warn("Failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if req.Group == "" || req.Song == "" {
		h.logger.Warn("Group and Song fields are required")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Group and Song fields are required",
			Message: "Please provide both group and song names",
		})
	}
	newSong, err := h.serv.AddNewSong(req.Group, req.Song)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"group": req.Group,
			"song":  req.Song,
			"error": err,
		}).Error("Error adding new song")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to add new song",
			Message: err.Error(),
		})
	}

	h.logger.WithFields(logrus.Fields{
		"group": newSong.Group,
		"song":  newSong.Song,
	}).Info("New song added successfully")

	return ctx.Status(fiber.StatusCreated).JSON(DataResponseSong{
		Data:    newSong,
		Message: "Song added successfully",
	})
}

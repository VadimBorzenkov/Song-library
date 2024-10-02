package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/sirupsen/logrus"
)

var possibleDateFormats = []string{
	"2006-01-02",
	"02-01-2006",
	"02/01/2006",
	"2006-01-02 15:04:05",
	"02-01-2006 15:04:05",
	"02/01/2006 15:04:05",
	"January 2, 2006",
	"02.01.2006",
}

func (h *ApiService) parseAndFormatDate(dateStr string) (string, error) {
	var parsedDate time.Time
	var err error

	for _, format := range possibleDateFormats {
		parsedDate, err = time.Parse(format, dateStr)
		if err == nil {
			formattedDate := parsedDate.Format(DateFormat)
			h.logger.WithFields(logrus.Fields{
				"inputDate":     dateStr,
				"parsedDate":    parsedDate,
				"formattedDate": formattedDate,
				"usedFormat":    format,
			}).Info("Successfully parsed and formatted date")
			return formattedDate, nil
		}
		h.logger.WithFields(logrus.Fields{
			"inputDate": dateStr,
			"format":    format,
			"error":     err,
		}).Debug("Failed to parse date with format")
	}

	h.logger.WithField("inputDate", dateStr).Error("Failed to parse date in all formats")
	return "", errors.New("invalid date format")
}

func (s *ApiService) GetSongsWithPaginate(filter map[string]string, limit, offset int) ([]models.Song, error) {
	songs, err := s.repo.GetData(filter, limit, offset)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"filter": filter,
			"limit":  limit,
			"offset": offset,
		}).Error("Failed to fetch songs: ", err)
		return nil, err
	}

	return songs, nil
}

func (s *ApiService) GetSongWithVerses(id, limit, offset int) (*models.Song, error) {
	song, err := s.repo.GetSongPagi(id, limit, offset)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"songID": id,
			"limit":  limit,
			"offset": offset,
		}).Error("Failed to fetch song with verses: ", err)
		return nil, err
	}

	s.logger.Infof("Successfully fetched song '%s' with %d verses", song.Song, limit)
	return song, nil
}

func (s *ApiService) AddNewSong(group, song string) (*models.Song, error) {
	songDetail, err := s.exApi.FetchSongInfo(group, song)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"group": group,
			"song":  song,
		}).Error("Failed to fetch song details from external API: ", err)
		return nil, err
	}

	formattedDate, err := s.parseAndFormatDate(songDetail.ReleaseDate)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"group":       group,
			"song":        song,
			"releaseDate": songDetail.ReleaseDate,
		}).Error("Failed to parse release date: ", err)
		return nil, err
	}

	newSong := &models.Song{
		Group:       group,
		Song:        song,
		ReleaseDate: formattedDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	err = s.repo.AddNewSong(newSong)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"group": newSong.Group,
			"song":  newSong.Song,
		}).Error("Failed to add new song to the database: ", err)
		return nil, err
	}

	return newSong, nil
}

func (s *ApiService) UpdateSong(song *models.Song) error {
	err := s.repo.UpdateSongData(song)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"songID": song.ID,
			"song":   song.Song,
		}).Error("Failed to update song: ", err)
		return err
	}

	return nil
}

func (s *ApiService) DeleteSong(id int) error {

	rowsAffected, err := s.repo.DeleteSong(id)
	if err != nil {
		s.logger.WithField("songID", id).Error("Failed to delete song: ", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("song with ID %d not found", id)
	}

	s.logger.Infof("Successfully deleted song with ID %d", id)
	return nil
}

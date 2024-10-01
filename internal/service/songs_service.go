package service

import (
	"errors"
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

	h.logger.WithField("inputDate", dateStr).Debug("Attempting to parse date")

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
	s.logger.WithFields(logrus.Fields{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
	}).Debug("Fetching songs with pagination")

	songs, err := s.repo.GetData(filter, limit, offset)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"filter": filter,
			"limit":  limit,
			"offset": offset,
		}).Error("Failed to fetch songs: ", err)
		return nil, err
	}

	s.logger.Infof("Successfully fetched %d songs", len(songs))
	return songs, nil
}

func (s *ApiService) GetSongWithVerses(id, limit, offset int) (*models.Song, error) {
	s.logger.WithFields(logrus.Fields{
		"songID": id,
		"limit":  limit,
		"offset": offset,
	}).Debug("Fetching song with verses")

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
	s.logger.WithFields(logrus.Fields{
		"group": group,
		"song":  song,
	}).Info("Adding new song")

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

	s.logger.WithFields(logrus.Fields{
		"group":       newSong.Group,
		"song":        newSong.Song,
		"releaseDate": newSong.ReleaseDate,
	}).Debug("Adding new song to the database")

	err = s.repo.AddNewSong(newSong)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"group": newSong.Group,
			"song":  newSong.Song,
		}).Error("Failed to add new song to the database: ", err)
		return nil, err
	}

	s.logger.Infof("Successfully added new song '%s' by group '%s'", newSong.Song, newSong.Group)
	return newSong, nil
}

func (s *ApiService) UpdateSong(song *models.Song) error {
	s.logger.WithFields(logrus.Fields{
		"songID": song.ID,
		"song":   song.Song,
		"group":  song.Group,
	}).Info("Updating song")

	formattedDate, err := s.parseAndFormatDate(song.ReleaseDate)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"songID":      song.ID,
			"releaseDate": song.ReleaseDate,
		}).Error("Failed to parse release date: ", err)
		return err
	}
	song.ReleaseDate = formattedDate

	err = s.repo.UpdateSongData(song)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"songID": song.ID,
			"song":   song.Song,
		}).Error("Failed to update song: ", err)
		return err
	}

	s.logger.Infof("Successfully updated song '%s' by group '%s'", song.Song, song.Group)
	return nil
}
func (s *ApiService) DeleteSong(id int) error {
	s.logger.WithField("songID", id).Info("Deleting song")

	err := s.repo.DeleteSong(id)
	if err != nil {
		s.logger.WithField("songID", id).Error("Failed to delete song: ", err)
		return err
	}

	s.logger.Infof("Successfully deleted song with ID %d", id)
	return nil
}

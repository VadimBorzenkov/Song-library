package service

import (
	"github.com/VadimBorzenkov/online-song-library/internal/models"
	"github.com/sirupsen/logrus"
)

// Получение песен с фильтрацией и пагинацией
func (s *ApiService) GetSongsWithPaginate(filter map[string]string, limit, offset int) ([]models.Song, error) {
	// Логируем входные параметры фильтрации
	s.logger.WithFields(logrus.Fields{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
	}).Debug("Fetching songs with pagination")

	// Получаем данные через репозиторий
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

// Получение текста песни с пагинацией по куплетам
func (s *ApiService) GetSongWithVerses(id, limit, offset int) (*models.Song, error) {
	// Логируем параметры запроса
	s.logger.WithFields(logrus.Fields{
		"songID": id,
		"limit":  limit,
		"offset": offset,
	}).Debug("Fetching song with verses")

	// Получаем песню с пагинацией по куплетам через репозиторий
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

// Добавление новой песни
func (s *ApiService) AddNewSong(group, song string) (*models.Song, error) {
	s.logger.WithFields(logrus.Fields{
		"group": group,
		"song":  song,
	}).Info("Adding new song")

	// Запрос во внешний API для получения данных о песне
	songDetail, err := s.exApi.FetchSongInfo(group, song)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"group": group,
			"song":  song,
		}).Error("Failed to fetch song details from external API: ", err)
		return nil, err
	}

	// Создаем новую запись для добавления в БД
	newSong := &models.Song{
		Group:       group,
		Song:        song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// Логируем данные новой песни
	s.logger.WithFields(logrus.Fields{
		"group":       newSong.Group,
		"song":        newSong.Song,
		"releaseDate": newSong.ReleaseDate,
	}).Debug("Adding new song to the database")

	// Добавляем песню в базу данных через репозиторий
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

// Обновление песни
func (s *ApiService) UpdateSong(song *models.Song) error {
	s.logger.WithFields(logrus.Fields{
		"songID": song.ID,
		"song":   song.Song,
		"group":  song.Group,
	}).Info("Updating song")

	// Обновляем данные песни через репозиторий
	err := s.repo.UpdateSongData(song)
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

// Удаление песни
func (s *ApiService) DeleteSong(id int) error {
	s.logger.WithField("songID", id).Info("Deleting song")

	// Удаляем песню через репозиторий
	err := s.repo.DeleteSong(id)
	if err != nil {
		s.logger.WithField("songID", id).Error("Failed to delete song: ", err)
		return err
	}

	s.logger.Infof("Successfully deleted song with ID %d", id)
	return nil
}

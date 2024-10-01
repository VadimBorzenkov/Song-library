package routes

import (
	"github.com/VadimBorzenkov/online-song-library/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
)

func RegistrationRoutes(app *fiber.App, h handler.Handler) {
	songsRoutes := app.Group("/songs")

	songsRoutes.Get("/", h.GetSongs)
	songsRoutes.Get("/get_song", h.GetSongWithVerses)
	songsRoutes.Post("/add_song", h.AddNewSong)
	songsRoutes.Put("/update_song", h.UpdateSong)
	songsRoutes.Delete("/:id", h.DeleteSong)
}

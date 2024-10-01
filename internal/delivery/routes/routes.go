package routes

import (
	"github.com/VadimBorzenkov/online-song-library/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func RegistrationRoutes(app *fiber.App, h handler.Handler) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	songsRoutes := app.Group("/songs")

	songsRoutes.Get("/", h.GetSongs)
	songsRoutes.Get("/get_song/:id", h.GetSongWithVerses)
	songsRoutes.Post("/add_song", h.AddNewSong)
	songsRoutes.Put("/update_song/:id", h.UpdateSong)
	songsRoutes.Delete("/delete_song/:id", h.DeleteSong)

	//Including swagger
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
}

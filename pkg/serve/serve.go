package serve

import (
	"github.com/gitu/katastasi/pkg/types"
	"github.com/gitu/katastasi/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"net/http"
	"time"
)

var lastUpdate = time.Now() // time.UnixMilli(0)

func StartServer(info string) {

	app := fiber.New()
	app.Use(compress.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/api/status/:id", func(c *fiber.Ctx) error {
		if c.Params("id") == "s1" {
			return c.JSON(types.StatusInfo{
				LastUpdate:    lastUpdate.UnixMilli(),
				StatusPage:    types.StatusPage{Id: "s1", Name: "Nikephoros"},
				OverallStatus: types.OK,
				Services: []types.Service{
					{
						Id:     "s1-1",
						Name:   "Nikephoros API",
						Status: types.OK,
					},
					{
						Id:     "s1-2",
						Name:   "Nikephoros UI",
						Status: types.OK,
					},
				},
			})
		} else if c.Params("id") == "s2" {
			return c.JSON(types.StatusInfo{
				LastUpdate:    lastUpdate.UnixMilli(),
				StatusPage:    types.StatusPage{Id: "s2", Name: "Okeanos"},
				OverallStatus: types.Warning,
				Services: []types.Service{
					{
						Id:     "s2-1",
						Name:   "Okeanos API",
						Status: types.Warning,
					},
					{
						Id:     "s2-2",
						Name:   "Okeanos UI",
						Status: types.OK,
					},
				},
			})
		}

		return c.SendStatus(http.StatusNotFound)
	})

	app.Get("/api/status", func(c *fiber.Ctx) error {
		return c.JSON([]types.StatusPage{
			{
				Id:   "s1",
				Name: "Nikephoros",
			},
			{
				Id:   "s2",
				Name: "Okeanos",
			},
			{
				Id:   "s3",
				Name: "Shine-Tsu-Hiko",
			},
		})
	})
	app.Get("/api/info", func(c *fiber.Ctx) error {
		return c.SendString(info)
	})

	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		return c.SendString("User-agent: *\nDisallow: /")
	})

	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(ui.DistDir),
		PathPrefix: "dist/assets",
		Browse:     false,
		MaxAge:     30 * 24 * 60 * 60,
	}))
	app.Use("/favicon.png", filesystem.New(filesystem.Config{
		Root:       http.FS(ui.DistDir),
		PathPrefix: "dist",
		Browse:     false,
		MaxAge:     30 * 24 * 60 * 60,
	}))
	app.Use("/favicon.ico", filesystem.New(filesystem.Config{
		Root:       http.FS(ui.DistDir),
		PathPrefix: "dist",
		Browse:     false,
		MaxAge:     30 * 24 * 60 * 60,
	}))

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(ui.DistDir),
		PathPrefix: "dist",
		Browse:     false,
		Index:      "index.html",
		MaxAge:     0,
	}))

	err := app.Listen(":1323")
	if err != nil {
		log.Fatal("error starting server", err)
		return
	}
}

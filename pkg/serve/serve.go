package serve

import (
	"github.com/gitu/katastasi/pkg/core"
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

func StartServer(k *core.Katastasi) {

	app := fiber.New()
	app.Use(compress.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/api/env/:env/status/:page", func(c *fiber.Ctx) error {
		env := c.Params("env")
		page := c.Params("page")
		if _, f := k.Config.Environments[env]; !f {
			return c.SendStatus(http.StatusNotFound)
		}
		if _, f := k.Config.Environments[env].StatusPages[page]; !f {
			return c.SendStatus(http.StatusNotFound)
		}

		status := k.GetPageStatus(env, page)

		ret := types.StatusInfo{
			LastUpdate:    status.LastUpdate.UnixMilli(),
			OverallStatus: types.MapStatus(status.Status),
			Services:      []types.Service{},
			Name:          k.Config.Environments[env].StatusPages[page].Name,
		}
		for key, service := range status.Services {
			sd := k.Config.GetService(env, key)
			newService := types.Service{
				Id:         sd.ID,
				Name:       sd.Name,
				Status:     types.MapStatus(service.Status),
				LastUpdate: service.LastUpdate.UnixMilli(),
				Components: make([]types.ServiceComponent, len(service.Components)),
			}
			for i, component := range service.Components {
				newService.Components[i] = types.ServiceComponent{
					Status: types.MapStatus(component.Status),
					Name:   component.Name,
					Info:   component.StatusString,
				}
			}
			ret.Services = append(ret.Services, newService)
		}

		return c.JSON(ret)
	})

	app.Get("/api/envs", func(c *fiber.Ctx) error {
		environments := []types.Environment{}
		for _, env := range k.Config.Environments {
			e := types.Environment{
				Id:          env.ID,
				Name:        env.ID,
				StatusPages: map[string]string{},
				Services:    map[string]string{},
			}
			for _, sp := range env.StatusPages {
				e.StatusPages[sp.ID] = sp.Name
			}
			for _, s := range env.Services {
				e.Services[s.ID] = s.Name
			}
			environments = append(environments, e)
		}
		return c.JSON(environments)
	})
	app.Get("/api/info", func(c *fiber.Ctx) error {
		return c.SendString(k.KataStatus.Info)
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

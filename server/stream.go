package server

import (
	"html/template"
	"webRTC/handlers"
	"webRTC/models"
	"webRTC/services"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
    Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.Templates.ExecuteTemplate(w, name, data)
}

func WebRTCRoutes(e *echo.Echo, config *models.ConfigST) {
	streamService := services.NewStreamService()
	configService := services.NewConfigService()
	

	e.Static("/static", "web/static")
	e.Renderer = &Template{
		Templates: template.Must(template.ParseGlob("web/templates/*.tmpl")),
	  }

	e.GET(
		"/",
		handlers.HTTPAPIServerIndex(config),
	)

	e.POST(
		"/stream/receiver/:uuid",
		handlers.HTTPAPIServerStreamWebRTC(streamService, configService, config),
	)

	e.GET(
		"/stream/codec/:uuid",
		handlers.HTTPAPIServerStreamCodec(
			configService,
			config,
		),
	)

	e.GET(
		"/stream/player/:uuid",
		handlers.HTTPAPIServerStreamPlayer(
			config,
		),
	)
}
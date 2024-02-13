package server

import (
	"log"
	"webRTC/models"

	"github.com/labstack/echo/v4"
)

func NewServer() *echo.Echo {
	return echo.New()
}

func RunServer(e *echo.Echo, config *models.ConfigST) {
	WebRTCRoutes(e, config)
	SnapshotRoutes(e)
	log.Fatal(e.Start(":8888"))
}

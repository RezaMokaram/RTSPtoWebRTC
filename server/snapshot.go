package server

import (
	"webRTC/handlers"
	"webRTC/services"

	"github.com/labstack/echo/v4"
)

func SnapshotRoutes(e *echo.Echo) {
	snapshotService := services.NewSnapshotService()

	e.GET(
		"/snapshot",
		handlers.GetSnapshot(snapshotService),
	)
}
package handlers

import (
	"net/http"

	"webRTC/models"
	"webRTC/models/snapshot"
	"webRTC/services"

	"github.com/labstack/echo/v4"
)

func GetSnapshot(
	snapshotService services.SnapshotService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		snapshotRequest := new(snapshot.SnapshotRequest)

		if err := c.Bind(&snapshotRequest); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", err.Error()))
		}

		if !snapshotRequest.IsValid() {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", "empty fields or invalid fields"))
		}

		img, err := snapshotService.GetSnapshot(*snapshotRequest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("server error", err.Error()))
		}

		return c.Blob(http.StatusOK, "image/jpeg", img)
	}
}
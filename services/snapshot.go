package services

import (
	"webRTC/models/snapshot"
	
	"gocv.io/x/gocv"
)

type SnapshotService interface {
	GetSnapshot(req snapshot.SnapshotRequest) ([]byte, error)
}

type snapshotService struct {
}

func NewSnapshotService() SnapshotService {
	return &snapshotService{
	}
}

func (s *snapshotService) GetSnapshot(
	req snapshot.SnapshotRequest,
) (
	[]byte,
	 error,
) {
	webcam, err := gocv.OpenVideoCaptureWithAPI(req.RTSP, gocv.VideoCaptureFFmpeg)
	if err != nil {
		return make([]byte, 0), err
	}
	current := gocv.NewMat()
	webcam.Read(&current)
	encode, err := gocv.IMEncode(gocv.JPEGFileExt, current)
	if err != nil {
		return make([]byte, 0), err
	}
	bts := encode.GetBytes()
	return bts, nil
}
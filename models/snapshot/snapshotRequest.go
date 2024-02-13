package snapshot

type SnapshotRequest struct {
	RTSP string `json:"rtsp,omitempty"`
}

func (req SnapshotRequest) IsValid() bool {
	if len(req.RTSP) <= 0 {
		return false
	}
	return true
}
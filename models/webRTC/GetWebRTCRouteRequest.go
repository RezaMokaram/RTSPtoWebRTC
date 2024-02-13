package webRTC

type GetWebRTCRouteRequest struct {
	CameraRTSP string `json:"camera_ip,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (req GetWebRTCRouteRequest) IsValid() bool {
	if req.CameraRTSP == "" || req.Username == "" || req.Password == "" {
		return false
	}
	return true
}
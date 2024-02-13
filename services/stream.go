package services

import (
	"errors"
	"log"
	"time"
	"webRTC/models"

	"github.com/deepch/vdk/format/rtspv2"
)

var (
	ErrorStreamExitNoVideoOnStream = errors.New("Stream Exit No Video On Stream")
	ErrorStreamExitRtspDisconnect  = errors.New("Stream Exit Rtsp Disconnect")
	ErrorStreamExitNoViewer        = errors.New("Stream Exit On Demand No Viewer")
)

type StreamService interface {
	ServeStreams(config *models.ConfigST)
	RTSPWorkerLoop(name, url string, OnDemand bool, config *models.ConfigST)
	RTSPWorker(name, url string, OnDemand bool, config *models.ConfigST) error
}

type streamService struct {
}

func NewStreamService() StreamService {
	return &streamService{
	}
}

func (s *streamService) ServeStreams(config *models.ConfigST) {
	// config := NewConfigService().LoadConfig()
	for k, v := range config.Streams {
		if !v.OnDemand {
			go s.RTSPWorkerLoop(k, v.URL, v.OnDemand, config)
		}
	}
}

func (s *streamService) RTSPWorkerLoop(name, url string, OnDemand bool, config *models.ConfigST) {
	// config := NewConfigService().LoadConfig()
	defer config.RunUnlock(name)
	for {
		log.Println("Stream Try Connect", name)
		err := s.RTSPWorker(name, url, OnDemand, config)
		if err != nil {
			log.Println(err)
		}
		if OnDemand && !config.HasViewer(name) {
			log.Println(ErrorStreamExitNoViewer)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *streamService) RTSPWorker(name, url string, OnDemand bool, config *models.ConfigST) error {
	// config := NewConfigService().LoadConfig()
	keyTest := time.NewTimer(20 * time.Second)
	clientTest := time.NewTimer(20 * time.Second)
	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: url, DisableAudio: false, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: false})
	if err != nil {
		return err
	}
	defer RTSPClient.Close()
	if RTSPClient.CodecData != nil {
		config.CoAd(name, RTSPClient.CodecData)
	}
	var AudioOnly bool
	if len(RTSPClient.CodecData) == 1 && RTSPClient.CodecData[0].Type().IsAudio() {
		AudioOnly = true
	}
	for {
		select {
		case <-clientTest.C:
			if OnDemand && !config.HasViewer(name) {
				return ErrorStreamExitNoViewer
			}
		case <-keyTest.C:
			return ErrorStreamExitNoVideoOnStream
		case signals := <-RTSPClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				config.CoAd(name, RTSPClient.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return ErrorStreamExitRtspDisconnect
			}
		case packetAV := <-RTSPClient.OutgoingPacketQueue:
			if AudioOnly || packetAV.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
			}
			config.Cast(name, *packetAV)
		}
	}
}

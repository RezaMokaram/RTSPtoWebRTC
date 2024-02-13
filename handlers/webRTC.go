package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"sort"

	"webRTC/models"
	"webRTC/services"

	"github.com/deepch/vdk/av"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"github.com/labstack/echo/v4"
)

//HTTPAPIServerIndex  index
func HTTPAPIServerIndex(
	config *models.ConfigST,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("-- HTTPAPIServerIndex --")
		_, all := config.List()
		if len(all) > 0 {
			c.Response().Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
			c.Redirect(http.StatusMovedPermanently, "stream/player/"+all[0])
		} else {
			// c.HTML(http.StatusOK, "index.tmpl", gin.H{
			// 	"port":    config.Server.HTTPPort,
			// 	"version": time.Now().String(),
			// })
			data := map[string]interface{}{
					"port":    config.Server.HTTPPort,
					"version": time.Now().String(),
			}
			return c.Render(http.StatusOK, "index.tmpl", data)
		}
		return c.Render(http.StatusOK, "index.tmpl", nil)
	}
}

// HTTPAPIServerStreamPlayer stream player
func HTTPAPIServerStreamPlayer(
	config *models.ConfigST,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("-- HTTPAPIServerStreamPlayer --")
		_, all := config.List()
		sort.Strings(all)
		data := map[string]interface{}{
			"port":     config.Server.HTTPPort,
			"suuid":    c.Param("uuid"),
			"suuidMap": all,
			"version":  time.Now().String(),
		}
		fmt.Println("REACHED HERE BEFOR RENDER...")
		return c.Render(http.StatusOK, "player.tmpl", data)
	}
}

//HTTPAPIServerStreamCodec stream codec
func HTTPAPIServerStreamCodec(
	configService services.ConfigService,
	config *models.ConfigST,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("-- IN ROTE HTTP API SERVER STREAM CODEC --")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		if config.Ext(c.Param("uuid")) {
			configService.RunIFNotRun(config, c.Param("uuid")) //todo
			codecs := config.CoGe(c.Param("uuid"))
			if codecs == nil {
				return c.JSON(http.StatusBadRequest, nil) // todo
			}
			var tmpCodec []models.JCodec
			for _, codec := range codecs {
				if codec.Type() != av.H264 && codec.Type() != av.PCM_ALAW && codec.Type() != av.PCM_MULAW && codec.Type() != av.OPUS {
					log.Println("Codec Not Supported WebRTC ignore this track", codec.Type())
					continue
				}
				if codec.Type().IsVideo() {
					tmpCodec = append(tmpCodec, models.JCodec{Type: "video"})
				} else {
					tmpCodec = append(tmpCodec, models.JCodec{Type: "audio"})
				}
			}
			fmt.Println("codec --->>>  ", tmpCodec)
			return c.JSON(http.StatusOK, tmpCodec)
		}
		fmt.Println("codec --->>>  empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
}

//HTTPAPIServerStreamWebRTC stream video over WebRTC
func HTTPAPIServerStreamWebRTC(
	webRTCService services.StreamService, 
	configService services.ConfigService, 
	config *models.ConfigST,
	) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("-- HTTPAPIServerStreamWebRTC --")
		if !config.Ext(c.FormValue("suuid")) {
			log.Println("Stream Not Found")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		configService.RunIFNotRun(config, c.FormValue("suuid")) // todo error
		codecs := config.CoGe(c.FormValue("suuid"))
		if codecs == nil {
			log.Println("Stream Codec Not Found")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		var AudioOnly bool
		if len(codecs) == 1 && codecs[0].Type().IsAudio() {
			AudioOnly = true
		}
		muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: config.GetICEServers(), PortMin: config.GetWebRTCPortMin(), PortMax: config.GetWebRTCPortMax()})
		answer, err := muxerWebRTC.WriteHeader(codecs, c.FormValue("data"))
		if err != nil {
			log.Println("WriteHeader", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		
		go func() {
			cid, ch := config.ClAd(c.FormValue("suuid"))
			defer config.ClDe(c.FormValue("suuid"), cid)
			defer muxerWebRTC.Close()
			var videoStart bool
			noVideo := time.NewTimer(10 * time.Second)
			for {
				select {
				case <-noVideo.C:
					log.Println("noVideo")
					return
				case pck := <-ch:
					if pck.IsKeyFrame || AudioOnly {
						noVideo.Reset(10 * time.Second)
						videoStart = true
					}
					if !videoStart && !AudioOnly {
						continue
					}
					err = muxerWebRTC.WritePacket(pck)
					if err != nil {
						log.Println("WritePacket", err)
						return
					}
				}
			}
		}()
		// _, err = c.Writer.Write([]byte(answer))
		// if err != nil {
		// 	log.Println("Write", err)
		// 	return c.JSON(http.StatusInternalServerError, nil)
		// }
		return c.String(http.StatusOK, answer)
	}
}
package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	_ "sync"

	"webRTC/models"
)

type ConfigService interface {
	LoadConfig() (*models.ConfigST)
	RunIFNotRun(element *models.ConfigST, uuid string)
}

type configService struct {
}

func NewConfigService() ConfigService {
	return &configService{
	}
}

func (s *configService) LoadConfig() *models.ConfigST {
	var tmp models.ConfigST
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	for i, v := range tmp.Streams {
		v.Cl = make(map[string]models.Viewer)
		tmp.Streams[i] = v
	}
	fmt.Println(tmp.Streams)
	return &tmp
}


func (s *configService) RunIFNotRun(element *models.ConfigST, uuid string) {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	if tmp, ok := element.Streams[uuid]; ok {
		if tmp.OnDemand && !tmp.RunLock {
			tmp.RunLock = true
			element.Streams[uuid] = tmp
			streamService := NewStreamService()
			go streamService.RTSPWorkerLoop(uuid, tmp.URL, tmp.OnDemand, element)
		}
	}
}
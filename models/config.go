package models

import (
	"sync"
	"crypto/rand"
	"fmt"
	"time"
	_ "sync"
	// "encoding/json"
	// "io/ioutil"
	// "log"

	"github.com/deepch/vdk/av"
)

//ConfigST struct
type ConfigST struct {
	Mutex   sync.RWMutex
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

//ServerST struct
type ServerST struct {
	HTTPPort      string   `json:"http_port"`
	ICEServers    []string `json:"ice_servers"`
	WebRTCPortMin uint16   `json:"webrtc_port_min"`
	WebRTCPortMax uint16   `json:"webrtc_port_max"`
}

//StreamST struct
type StreamST struct {
	URL      string `json:"url"`
	Status   bool   `json:"status"`
	OnDemand bool   `json:"on_demand"`
	RunLock  bool   `json:"-"`
	Codecs   []av.CodecData
	Cl       map[string]Viewer
}

type Viewer struct {
	C chan av.Packet
}

type JCodec struct {
	Type string
}

// Methods

func (element *ConfigST) RunUnlock(uuid string) {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	if tmp, ok := element.Streams[uuid]; ok {
		if tmp.OnDemand && tmp.RunLock {
			tmp.RunLock = false
			element.Streams[uuid] = tmp
		}
	}
}

func (element *ConfigST) HasViewer(uuid string) bool {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	if tmp, ok := element.Streams[uuid]; ok && len(tmp.Cl) > 0 {
		return true
	}
	return false
}

func (element *ConfigST) GetICEServers() []string {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	return element.Server.ICEServers
}

func (element *ConfigST) GetWebRTCPortMin() uint16 {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	return element.Server.WebRTCPortMin
}

func (element *ConfigST) GetWebRTCPortMax() uint16 {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	return element.Server.WebRTCPortMax
}

func (element *ConfigST) Cast(uuid string, pck av.Packet) {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	for _, v := range element.Streams[uuid].Cl {
		if len(v.C) < cap(v.C) {
			v.C <- pck
		}
	}
}

func (element *ConfigST) Ext(suuid string) bool {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	_, ok := element.Streams[suuid]
	return ok
}

func (element *ConfigST) CoAd(suuid string, codecs []av.CodecData) StreamST  {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	t := element.Streams[suuid]
	t.Codecs = codecs
	element.Streams[suuid] = t
	fmt.Println("in another wiered func --> ", t.Codecs)
	return t
}

func (element *ConfigST) CoGe(suuid string) []av.CodecData {
	for i := 0; i < 100; i++ {
		element.Mutex.RLock()
		
		tmp, ok := element.Streams[suuid]
		fmt.Println("in wiered func search for this:  ", suuid)
		element.Mutex.RUnlock()
		if !ok {
			fmt.Println("in wiered : return nil  ")
			return nil
		}
		if tmp.Codecs != nil {
			fmt.Println("in wiered : return ok: ", tmp.Codecs)
			return tmp.Codecs
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

func (element *ConfigST) ClAd(suuid string) (string, chan av.Packet) {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	cuuid := PseudoUUID()
	ch := make(chan av.Packet, 100)
	element.Streams[suuid].Cl[cuuid] = Viewer{C: ch}
	return cuuid, ch
}

func (element *ConfigST) List() (string, []string) {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	var res []string
	var fist string
	for k := range element.Streams {
		if fist == "" {
			fist = k
		}
		res = append(res, k)
	}
	return fist, res
}

func (element *ConfigST) ClDe(suuid, cuuid string) {
	element.Mutex.Lock()
	defer element.Mutex.Unlock()
	delete(element.Streams[suuid].Cl, cuuid)
}

func PseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
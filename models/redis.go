package models

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bayugyug/rest-api-throttleip/utils"
	"github.com/google/uuid"
	redis "gopkg.in/redis.v3"
)

const (
	IPDeniedKey  = "THROTTLE::IP::DENIED"
	IPAllowedKey = "THROTTLE::IP::ALLOWED"
)

type TrackerIPHistory struct {
	lock           sync.Mutex
	HistoryChannel chan *TrackerIP
}

func NewTrackerIPHistory() *TrackerIPHistory {
	return &TrackerIPHistory{
		HistoryChannel: make(chan *TrackerIP, 5000),
	}
}

var IPHistoryLogs map[string]int

//ManageQ the ip history logs
func (h *TrackerIPHistory) ManageQ(isReady chan bool) {

	//get new set of history
	IPHistoryLogs = h.InitQ()

	//now its minute ;-)
	ticker := time.NewTicker(time.Second * 60)

	//ready
	isReady <- true
	utils.Dumper("ManageQ::IsReady")
	for {
		select {
		case <-ticker.C:
			//init every n minute
			IPHistoryLogs = h.InitQ()
			utils.Dumper("history::q refresh")
		}
	}
}

//InitQ set new map of queue data
func (h *TrackerIPHistory) InitQ() map[string]int {
	return make(map[string]int)
}

//GetIP return count
func (h *TrackerIPHistory) GetIP(s string) int {
	//just in case ;-)
	h.lock.Lock()
	defer h.lock.Unlock()
	if _, oks := IPHistoryLogs[s]; !oks {
		IPHistoryLogs[s] = 1
	} else {
		IPHistoryLogs[s]++
	}
	utils.Dumper("history::q", s, IPHistoryLogs[s])
	//give it back
	return IPHistoryLogs[s]
}

//ManageHistory
func (h *TrackerIPHistory) ManageHistory(isReady chan bool, cache *redis.Client) {

	//ready
	isReady <- true
	utils.Dumper("ManageHistory::IsReady")
	pipe := cache.Pipeline()
	for {
		select {
		case info := <-h.HistoryChannel:
			if info.IP != "" {
				data, derr := json.Marshal(info)
				if derr != nil {
					log.Println("FAILED_TO_ADD_REDIS", derr)
					continue
				}
				key := IPAllowedKey
				if strings.EqualFold(info.Status, "Denied") {
					key = IPDeniedKey
				}
				pipe.HSet(key, time.Now().Format("20060102-150405")+"::"+uuid.New().String()+"::"+info.IP, string(data)) //no expiry on the summary list
				_, err := pipe.Exec()
				if err != nil {
					log.Println("FAILED_TO_ADD_REDIS", err)
				}

			}
		}
	}
}

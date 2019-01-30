package utils

import (
	"encoding/json"
	"log"
)

var ShowMeLog = true

func Dumper(infos ...interface{}) {
	if ShowMeLog {
		j, _ := json.MarshalIndent(infos, "", "\t")
		log.Println(string(j))
	}
}

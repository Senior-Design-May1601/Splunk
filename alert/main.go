package alert

import (
	"encoding/json"
	"log"
)

type Alert struct {
	Source     string            `json:"source"`
	SourceType string            `json:"sourcetype"`
	Event      map[string]string `json:"event"`
}

func NewSplunkAlertMessage(meta map[string]string) string {
	msg, err := json.Marshal(Alert{Source: "projectmain",
		SourceType: "alert",
		Event:      meta})
	if err != nil {
		log.Fatal("json marshal error:", err)
	}
	return string(msg)
}

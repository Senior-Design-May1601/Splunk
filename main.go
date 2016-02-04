package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Senior-Design-May1601/Splunk/alert"
	"github.com/Senior-Design-May1601/projectmain/loggerplugin"
)

type SplunkPlugin struct {
	Client *http.Client
	Token  string
	Url    string
}

type GenericAlert struct {
	Event []byte `json:"event"`
}

func (s *SplunkPlugin) Log(msg []byte, _ *int) error {
	var logMsg []byte
	var splunkAlert alert.Message

	// check if we got a well formed splunk alert. if not, build a generic alert
	err := json.Unmarshal(msg, &splunkAlert)
	if err != nil {
		// didn't get a splunk alert. just pack into a generic alert
		logMsg, err = json.Marshal(GenericAlert{msg})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		logMsg, err = json.Marshal(splunkAlert)
		if err != nil {
			log.Fatal(err)
		}
	}

	req, err := http.NewRequest("POST", s.Url, bytes.NewBuffer(logMsg))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Splunk %s", s.Token))

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(body))
	}

	return nil
}

var client *http.Client

func main() {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}

	client = &http.Client{Transport: tr}
	plugin := &SplunkPlugin{Client: client}

	p, err := loggerplugin.NewLoggerPlugin(plugin)
	if err != nil {
		log.Fatal(err)
	}

	err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

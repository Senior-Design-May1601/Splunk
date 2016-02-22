package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/Senior-Design-May1601/Splunk/alert"
	"github.com/Senior-Design-May1601/projectmain/loggerplugin"
)

type Config struct {
	Endpoint []endpoint
}

type endpoint struct {
	Host      string
	Port      int
	URL       string
	AuthToken string
	RootCAs   []string
}

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

	// check if we got a well formed splunk alert.
	// if not, build a generic alert
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
	configPath := flag.String("config", "", "config file")
	flag.Parse()

	var configs Config
	if _, err := toml.DecodeFile(*configPath, &configs); err != nil {
		log.Fatal(err)
	}
	// TODO: consider supporting multiple endpoints. for now just always
	//       use the first one
	config := configs.Endpoint[0]

	var roots *x509.CertPool

	roots = nil

	if len(config.RootCAs) > 0 {
		roots = x509.NewCertPool()
		for _, CA := range config.RootCAs {
			pem, err := ioutil.ReadFile(CA)
			if err != nil {
				log.Fatal(err)
			}
			ok := roots.AppendCertsFromPEM(pem)
			if !ok {
				log.Fatal("failed to parse CA certificate")
			}
		}
	}

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: roots},
		DisableCompression: true,
	}

	client = &http.Client{Transport: tr}
	plugin := &SplunkPlugin{Client: client,
		Token: config.AuthToken,
		Url:   "https://" + config.Host + ":" + strconv.Itoa(config.Port) + config.URL}

	p, err := loggerplugin.NewLoggerPlugin(plugin)
	if err != nil {
		log.Fatal(err)
	}

	err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

package main
import(
//	"github.com/Senior-Design-May1601/projectmain/loggerplugin"
	"bytes"
	"net/http"
	"fmt"
	"encoding/json"
	"crypto/tls"
	"github.com/BurntSushi/toml"
	"log"
	"errors"
	"os"
	"time"
)

type SplunkPlugin struct{
	Token string
	Url string
}

type Event struct{
	Msg  interface{} `json:"event"`
}
type SplunkConfig struct{
	Token string
	Url  string
}

type Cache struct{
	Time time.Time
	Event []byte
}

func (s *SplunkPlugin)EncodeMSG(msg string, _*int) error{
	e := Event{Msg : msg}
	data := map[string]Event{"event":Event{Msg:msg}}
	encoded,err := json.Marshal(data)

	if err != nil {
		fmt.Println(e)
	}
	var reply *int
	s.Log(encoded ,reply)
	return nil 

}

func (s *SplunkPlugin) Log(msg []byte, _*int) error{
	transport := &http.Transport{TLSClientConfig:&tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport : transport}
	
	reader := bytes.NewReader(msg)
	req,_ := http.NewRequest("POST",s.Url,reader)
	req.Header.Add("Authorization",fmt.Sprintf("Splunk %s",s.Token))
	resp, err :=client.Do(req)

	if err != nil{
		log.Fatal(err)
	}
	if resp.StatusCode != 200{
		return errors.New(resp.Status)
	}

	return nil
}

func CacheEvent(msg []byte) {
	f, err := os.OpenFile("requests.toml",os.O_APPEND|os.O_WRONLY,0600)
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()
	
	m := Cache{Time: time.Now().UTC(),Event:msg}
	c := make(map[string]Cache)
	c["0"] = m
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(c); err != nil{
		log.Fatal(err)
	}	
}

func ReadCache() []Cache{
	var events []Cache
	if _,err :=toml.DecodeFile("requests.toml",&events); err != nil{
		log.Fatal(err)
	}
	return events
}

func main(){
	var config SplunkConfig
	if _, err := toml.DecodeFile("config.toml",&config); err != nil {
		log.Fatal(err)
	}
/*	p, err := loggerplugin.NewLoggerPlugin(&SplunkPlugin{
								Token: config.Token,
								Url : config.Url,
								})

	if err != nil {
		log.Fatal(err)
	}
	err = p.Run()
	if err != nil{
		log.Fatal(err)
	}*/

	s := SplunkPlugin{Token:config.Token, Url :config.Url,}
	var reply *int
	s.EncodeMSG("Another Alert",reply)

	data, _ := json.Marshal("hello")
	CacheEvent(data)	
}

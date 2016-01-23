package main
import(
//	"github.com/Senior-Design-May1601/projectmain/loggerplugin"
//	"log"
	"bytes"
	"net/http"
	"fmt"
	"encoding/json"
	"crypto/tls"
	"github.com/BurntSushi/toml"
	"log"
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
		fmt.Println(err)
	}else{
		fmt.Println(resp)
	}
	return nil
}

func main(){
/*	p, err := loggerplugin.NewLoggerPlugin(&SplunkPlugin{
								Token: "1280B3F3-3AEC-49B5-861D-49E745BFB827",
								Url : "https://localhost:8088/services/collector",
								})

	if err != nil {
		log.Fatal(err)
	}
	err = p.Run()
	if err != nil{
		log.Fatal(err)
	}*/
	var config SplunkConfig
	if _, err := toml.DecodeFile("config.toml",&config); err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.Token)
	s := SplunkPlugin{Token:config.Token, Url :config.Url,}
	var reply *int
	s.EncodeMSG("Another Alert",reply)
}

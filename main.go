package main
import(
	"github.com/Senior-Design-May1601/projectmain/loggerplugin"
        "net/http"
        "github.com/BurntSushi/toml"
        "log"
        "fmt"
        "bytes"
        "crypto/tls"
	"encoding/json"
	"os"
	"errors"
	"io/ioutil"
)

type splunkConfig struct{
        Token string
        Url string
}

type SplunkPlugin struct{
        SplunkConfig []splunkConfig
        Client *http.Client
}

type SplunkAlert struct{
        Service string `json:"source"`
        Meta    map[string]string `json:"event"`
}
type GenericAlert struct{
        Event interface{} `json:"event"`
}


func (s *SplunkPlugin) Log(msg []byte, _*int) error{
        var t SplunkAlert
        var g GenericAlert
        if err := json.Unmarshal(msg,&t); err != nil{
               	if err := json.Unmarshal(msg,&g.Event); err != nil{
			g.Event=string(msg)	
		}
                msg,err = json.Marshal(g)
                if err != nil {	
                        return err
                }
        }
        request, err := http.NewRequest("POST",s.SplunkConfig[0].Url,bytes.NewBuffer(msg))
	
	if err != nil {
		return err
	}
        request.Header.Add("Authorization",fmt.Sprintf("Splunk %s",s.SplunkConfig[0].Token))
        resp,err :=s.Client.Do(request)

        if err != nil{
                return err
        }
	if resp.StatusCode != 200 {	
		body,_ := ioutil.ReadAll(resp.Body)
		return	errors.New(string(body))
	}

        return nil
}

func NewAlert() *SplunkAlert{
	
	return &SplunkAlert{}
}


func main(){
	
	tr := &http.Transport{
                TLSClientConfig:&tls.Config{InsecureSkipVerify : true},
        }
	
	client := &http.Client{Transport : tr}
	plugin := SplunkPlugin{Client:client}
	if _,err := toml.DecodeFile("dev-config.toml",&plugin);err != nil{
		log.Fatal(err)
	}
	f, err := os.Create("logfile")
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()
	for _, s := range plugin.SplunkConfig{
		f.WriteString(s.Token)
	}


	p, err := loggerplugin.NewLoggerPlugin(&plugin)
	if err != nil{
		log.Fatal(err)
	}	

	err = p.Run()
	if err != nil{
		log.Fatal(err)
	}
}

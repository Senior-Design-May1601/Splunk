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
)

type SplunkAlert struct{
        Service string `json:"source"`
        Meta    map[string]string `json:"event"`
}
type GenericAlert struct{
        Event interface{} `json:"event"`
}


type SplunkPlugin struct{
	Token string
	Url string
	Client *http.Client
}

func (s *SplunkPlugin) Log(msg []byte, _*int) error{
        var t SplunkAlert
        var g GenericAlert
        if err := json.Unmarshal(msg,&t); err != nil{
                json.Unmarshal(msg,&g.Event)
                msg,err = json.Marshal(g)
                if err != nil {
			log.Fatal(err)
                        return err
                }
        }
        request, err := http.NewRequest("POST",s.Url,bytes.NewBuffer(msg))
        request.Header.Add("Authorization",fmt.Sprintf("Splunk %s",s.Token))
        _,err =s.Client.Do(request)
        if err != nil{
		log.Fatal(err)
                return err
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
	if _,err := toml.DecodeFile("config.toml",&plugin);err != nil{
		log.Fatal(err)
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

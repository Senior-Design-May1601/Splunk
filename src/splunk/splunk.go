package splunk
import ("net/http"
		"fmt"
		"io"
		"io/ioutil"
		"bytes"
		"crypto/tls"
		"encoding/json"
		"os"
		"time"
		"bufio"
)

const(
 TOKEN = "5B6AB11C-FB81-4B06-A300-99E31FE3E781"  // Edit to client provided token from Splunk Event Collector
 EVENT = "https://localhost:8088/services/collector" // Change to Splunk server IP/services/collector
)
type Http struct{
	Method string `json:"method"` 
	Path string `json:"path"`
	Parameters map[string]string `json:"parameters"`
}

type Https struct{
	method string
	path string
	parameters map[string]string 
}

type SSH struct{
	username string
	password string
	key	string
}
type Event struct{
	SourceAddress string `json:"source address"`
	SourcePort int `json:"source port"`
	ServiceType string `json:"service type"`
	SSH *SSH `json:"ssh,omitempty"`
	Http *Http `json:"http,omitempty"`
	Https *Https `json:"https,omitempty"`
}


// HttpClient with TLS 
func httpClient() *http.Client{
	  tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr}
        return client
}

// Do five attempts, one every 5 seconds. If it still fails write to file
func cacheEvent(client *http.Client,request *http.Request,data []byte){

	f, err := os.OpenFile("cache.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(data); 
	f.WriteString("\n"); 
	
	
	f.Close()
	time.Sleep(5*time.Second)
	retryCache()
}
func retryCache(){
	f, err := os.OpenFile("cache.txt",os.O_RDONLY,0666)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	var line string
	var write []string
	scnr := bufio.NewScanner(f)
	
	for scnr.Scan(){
		line = scnr.Text()
		_,err:=request(line)
		if err!=nil{
			write = append(write,line)
			write = append(write,"\n")
			fmt.Println(write)
		}
	}
	if len(write)==0{
		fmt.Println("empty")
		os.Remove("cache.txt")
	}else{
		f,err = os.OpenFile("cache.txt",os.O_WRONLY,0666)
		for i:=0;i<len(write);i++{
			f.Write([]byte(write[i]))
		}
	}	
}
/*
* Usage example
*	m := map[string]string{"username":"user","password":"pass"}
*	h:=Http{Method:"POST",Path:"index.html",Parameters:m}
*	conn := Event{SourceAddress: "source",SourcePort: 239,ServiceType: "http",Http: &h,SSH : nil,Https:nil}
*	b,err:=conn.Send()
*/

func (e Event) Send() (string,error){
	client := httpClient()
	var payload io.Reader
    data:= map[string] Event{"event":e}  
	json_data,err := json.Marshal(data)
	
	//os.Stdout.Write(json_data)
	payload = bytes.NewReader(json_data)
	request, err:= http.NewRequest("POST", EVENT, payload)
	
	
	request.Header.Add("Authorization", fmt.Sprintf("Splunk %s", TOKEN))
	resp,err:=client.Do(request)
	if(err!=nil){
		go cacheEvent(client,request,json_data)
		return "",err
	}else{
		response,err:= ioutil.ReadAll(resp.Body)
		return string(response),err
	}
}
func request(e string) (*http.Response,error){
	client := httpClient()
	var payload io.Reader
    data:= map[string] string{"event":e}  
	json_data,err := json.Marshal(data)
	payload = bytes.NewReader(json_data)
	request, err:= http.NewRequest("POST", EVENT, payload)
	request.Header.Add("Authorization", fmt.Sprintf("Splunk %s", TOKEN))
	resp,err:=client.Do(request)
	
	return resp, err
}
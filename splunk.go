package main
import ("net/http"
		"fmt"
		"io"
		"io/ioutil"
		"bytes"
		"crypto/tls"
		"encoding/json"
)

const(
 TOKEN = "5B6AB11C-FB81-4B06-A300-99E31FE3E781"  // Edit to client provided token from Splunk Event Collector
 EVENT = "https://localhost:8088/services/collector" // Change to Splunk server IP/services/collector
)

// Throw error if exception found
func check(e error){
	if e!=nil{
		panic(e)
	}
}

// HttpClient with TLS 
func httpClient() *http.Client{
	  tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr}
        return client
}

// Add event from string
func addEvent(_data string) (string,error){
	client := httpClient()
	var payload io.Reader
     
	data := map[string]string{"event":_data} 
	json_data,err := json.Marshal(data)
	check(err)
	
	payload = bytes.NewReader(json_data)
	request, err:= http.NewRequest("POST", EVENT, payload)
	check(err)
	
	request.Header.Add("Authorization", fmt.Sprintf("Splunk %s", TOKEN))
	resp,err:=client.Do(request)
	response,err:= ioutil.ReadAll(resp.Body)
	
	return string(response),err
}

//Create event from a file
func addEventFromFile(filePath string,sourceType string) (string,error){
	_file, err := ioutil.ReadFile(filePath)
	check(err)
	content:=string(_file)
	client := httpClient()
	var data map[string]string
	
	if len(sourceType)>0{
		data = map[string]string{"event":content,"sourcetype":sourceType}
	}else{
		data = map[string]string{"event":content}
	}
	
	json_data,err := json.Marshal(data)
	check(err)
	
	payload := bytes.NewReader(json_data)
	request, err:= http.NewRequest("POST", EVENT, payload)
	check(err)
	
	request.Header.Add("Authorization", fmt.Sprintf("Splunk %s", TOKEN))
	resp,err:=client.Do(request)
	response,err:= ioutil.ReadAll(resp.Body)
	
	return string(response),err
}
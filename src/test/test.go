package main

import "splunk"
import "fmt"
import "time"

func main(){

	conn := splunk.Event{"192",80,"http",nil,nil,nil}
	b, err := conn.Send()
	if(err !=nil){
		fmt.Println(err)
	}else{
		fmt.Println(b)
	}
	
	
	time.Sleep(50 * time.Second)
}
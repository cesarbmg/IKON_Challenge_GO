package main

import (
	"bufio"
	"fmt"
	"os"
	"context"
	"log"
	"io/ioutil"
	"bytes"
  	"encoding/json"
    "net/http"
	"google.golang.org/grpc"	
	"src/github.com/cesarbmg/IKON_Challenge_GO/gRPC"
)

type device struct {  
    Capacity string `json:"Capacity"`
    Foreground string `json:"Foreground"`
    Background string `json:"Background"`
}

func gRPC(d device){
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:8083", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	client := Protocol.NewDeviceServiceClient(cc)
	request := &Protocol.DeviceRequest{Capacity: d.Capacity, Foreground: d.Foreground, Background: d.Background }
	resp, _ := client.Device(context.Background(), request)
	fmt.Println(resp.Response)
}

func rEST(d device){
	url:="http://localhost:8084"

	var jsonData []byte
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
	resp, err := client.Do(req)
	
    if err != nil {
		log.Fatal(err)
    }else{
    	defer resp.Body.Close()
		//debug
    	//fmt.Println("response Status:", resp.Status)
		data, _ := ioutil.ReadAll(resp.Body)		
		fmt.Println(string(data))
	}
}

func main() {
	readFile, err := os.Open("challenge.in")
 
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
 
	fileScanner := bufio.NewScanner(readFile)	
	fileScanner.Split(bufio.ScanLines)

	var fileTextLines []string
 
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
 
	readFile.Close()

	fmt.Println("INPUT...")

	for _, eachline := range fileTextLines {
		fmt.Println(eachline)
	}
	
	fmt.Println("OUTPUT...")
	var k = len(fileTextLines)

	for i := 0; i < k; i++ {
		var d device
		d.Capacity = fileTextLines[i]
		i++
		d.Foreground = fileTextLines[i]
		i++
		d.Background = fileTextLines[i]		
		
		gRPC(d)
		rEST(d)
	}
}

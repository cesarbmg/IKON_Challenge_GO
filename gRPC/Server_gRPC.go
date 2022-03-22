package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/cesarbmg/IKON_Challenge_GO/gRPC/Protocol"
	"google.golang.org/grpc"
)

var sURLgRPC string = "0.0.0.0"
var sPortgRPC string = "8084"
var addressgRPC string = sURLgRPC + ":" + sPortgRPC

type devicegRPC struct {
	Capacity   string `json:"Capacity"`
	Foreground string `json:"Foreground"`
	Background string `json:"Background"`
}

func rEST(d devicegRPC) string {
	url := "http://localhost:8083"

	var jsonData []byte
	jsonData, err := json.Marshal(d)

	if err != nil {
		log.Fatalf("error in convert json client REST: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("error in SendPOST client REST: %s", err)
	} else {
		defer resp.Body.Close()

		//debug
		//fmt.Println("response Status:", resp.Status)

		data, _ := ioutil.ReadAll(resp.Body)

		//debug
		//fmt.Println(string(data))

		return string(data)
	}

	return ""
}

//Estrcutura del server
type server struct {
}

func (*server) Device(ctx context.Context, request *Protocol.DeviceRequest) (*Protocol.DeviceResponse, error) {

	var d devicegRPC
	d.Capacity = request.Capacity
	d.Foreground = request.Foreground
	d.Background = request.Background

	var sResponse = rEST(d)

	fmt.Println(sResponse)

	response := &Protocol.DeviceResponse{
		Response: "gRPC => " + sResponse + "",
	}

	fmt.Println("Server gRPC is listening on " + addressgRPC + "...")

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", addressgRPC)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Println("Server gRPC is listening on " + addressgRPC + "...")
	s := grpc.NewServer()
	Protocol.RegisterDeviceServiceServer(s, &server{})
	s.Serve(lis)
}

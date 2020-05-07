package main

import (
	"bufio"
	"fmt"
	"os"
	"context"
	"log"
	"google.golang.org/grpc"	
	"github.com/cesarbmg/IKON_Challenge_GO/gRPC"
)

type device struct {  
    Capacity string `json:"Capacity"`
    Foreground string `json:"Foreground"`
    Background string `json:"Background"`
}

func gRPC(d device) string{
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:8084", opts)

	if err != nil {
		log.Fatalf("error in open gRPC: %s", err)
	}

	defer cc.Close()

	client := Protocol.NewDeviceServiceClient(cc)
	request := &Protocol.DeviceRequest{Capacity: d.Capacity, Foreground: d.Foreground, Background: d.Background }
	resp, _ := client.Device(context.Background(), request)

	fmt.Println(resp.Response)

	return resp.Response;
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

	var sResult []string

	var k = len(fileTextLines)
	for i := 0; i < k; i++ {
		var d device
		d.Capacity = fileTextLines[i]
		i++
		d.Foreground = fileTextLines[i]
		i++
		d.Background = fileTextLines[i]		
		
		sResult = append(sResult, gRPC(d))
		//sResult = append(sResult, rEST(d))
	}

	f, err1 := os.Create("challenge.out")
    if err1 != nil {
        log.Fatalf("Unable to create file output: %v", err1)
        f.Close()
        return
    }

	for _, s := range sResult {
        fmt.Fprintln(f, s)
	}
	
    err = f.Close()
    if err != nil {
		log.Fatalf("Unable to close file output: %v", err)
        return
	}
	
}

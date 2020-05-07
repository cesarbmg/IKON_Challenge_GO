package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"strconv"
	"google.golang.org/grpc"
	// "src/github.com/IKON_Challenge/gRPC"
	"Services"
)

//Estructura para los arreglos de tares Foreground y Background
type task struct{
	TaskID int
	Resource int
}

//Estructura para los arreglos de combinaciones de Id de Foreground y Background
type manager struct{
	TaskIDForeground int
	TaskIDBackground int
	ResourceTotal int
}

//Funciones para convertir string en int
func getInt(sInt string) int{
	sInt = strings.Trim(sInt," ")
	i, err := strconv.Atoi(sInt)
	if err == nil {
		return i;
	}
	return 0;
}

//Funcion para obtener los array de task desde string
func getArrayTasks(sTask string) []task{

	var tasks = []task {}

	re := regexp.MustCompile(`\((.*?)\)`)

	match := re.FindAllString(sTask, -1)
	for _, element := range match {
		element = strings.Trim(element, "(")
		element = strings.Trim(element, ")")

		//debug
		//fmt.Println(strings.Split(element, ",")[0]);
		//fmt.Println(strings.Split(element, ",")[1]);

		var I = getInt(strings.Split(element, ",")[0])
		var R = getInt(strings.Split(element, ",")[1])
		
		//debug
		//fmt.Println(I);
		//fmt.Println(R);

		item := task{
			TaskID: I, 
			Resource: R,
		}		
        tasks = append(tasks, item)
		
		//debug
		//fmt.Println(tasks);
	}

	//debug
	//fmt.Println(tasks);
	return tasks
}

//Funciones para obtener la maxima capacidad entre foreground y background
func getMaxManager(tasks []manager, capacity int) string{
	MaxCapacity := 0
	//var tasksMax = []manager {}
	sResponse :=  ""

	for _, task := range tasks {
		if  task.ResourceTotal == capacity{
			MaxCapacity = capacity
			break 
		}
		if MaxCapacity < task.ResourceTotal {
			MaxCapacity = task.ResourceTotal
		}
	}

	for _, task := range tasks {
		if  MaxCapacity == task.ResourceTotal{
			sResponse += strings.Join([]string {"(", strconv.Itoa(task.TaskIDForeground), ",", strconv.Itoa(task.TaskIDBackground), ")"}, "")
			//tasksMax = append(tasksMax, task)
		}
	}

	return sResponse

	//return tasksMax
}

//Estrcutura del server
type server struct {

}

func (*server) Device(ctx context.Context, request *Protocol.DeviceRequest) (*Protocol.DeviceResponse, error) {

	foreground:= getArrayTasks(request.Foreground)
	//debug
	//fmt.Println(foreground)

	background:= getArrayTasks(request.Background)
	//debug
	//fmt.Println(background)
	
	capacity:= getInt(request.Capacity)

	var tasks = []manager {}

	for _, taskf := range foreground {
		for _, taskb := range background {
			if (taskf.Resource + taskb.Resource) <= capacity {
				item := manager{
					TaskIDForeground: taskf.TaskID, 
					TaskIDBackground: taskb.TaskID, 
					ResourceTotal: (taskf.Resource + taskb.Resource),
				}
				tasks = append(tasks, item)
			}
		}
	}

	var sResponse = getMaxManager(tasks, capacity)

	fmt.Println(sResponse)

	response := &Protocol.DeviceResponse{
		Response: "gRPC => " + sResponse + "",
	}
	return response, nil
}

func main() {
	address :=  "0.0.0.0:8083"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error %v", err)
	}	
	fmt.Println("Server gRPC is listening on " + address + "...")
	s := grpc.NewServer()	
	Protocol.RegisterDeviceServiceServer(s, &server{})
	s.Serve(lis)
}
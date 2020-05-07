package main

import (
	"fmt"
	"log"
    "net/http"
	"regexp"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
)

type device struct {  
    Capacity string `json:"Capacity"`
    Foreground string `json:"Foreground"`
    Background string `json:"Background"`
}

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

type server struct{

}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//debug
	//fmt.Println(r.Method)	
	//fmt.Println(r.Body)

	switch r.Method {
    case "POST":
        w.WriteHeader(http.StatusOK)
   
		var d device

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		err1 := json.Unmarshal(body, &d)
		if err1 != nil {
			log.Println(err1)
		}
		
		//debug
		//fmt.Println(d)

		foreground:= getArrayTasks(d.Foreground)
		//debug
		//fmt.Println(foreground)

		background:= getArrayTasks(d.Background)
		//debug
		//fmt.Println(background)
		
		capacity:= getInt(d.Capacity)

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
		
		w.Write([]byte("REST => " + sResponse))
		
	default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
	}
}

func main() {
	address := "0.0.0.0:8084"
	s := &server{}	
	http.Handle("/", s)	
	fmt.Println("Server REST is listening on " + address + "...")	
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
}
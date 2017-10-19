/*basic implementation of client*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	flagc "github.com/fmarmol/flag"
	"github.com/fmarmol/todo/api"
	"google.golang.org/grpc"
)

const (
	LIST       string = "list"
	CREATE            = "create"
	UPDATE            = "update"
	DELETE            = "delete"
	LOW               = "LOW"
	MEDIUM            = "MEDIUM"
	HIGH              = "HIGH"
	TODO              = "TODO"
	INPROGRESS        = "INPROGRESS"
	DONE              = "DONE"
)

func main() {
	var startDate, endDate flagc.DateFlag

	host := flag.String("host", "localhost:8080", "host server")
	action := flag.String("action", LIST, "choose action from list|create|update|delete")
	uid := flag.Int("uid", 1, "uuid task")
	subject := flag.String("subject", "another task", "subject of the task")
	priority := flag.String("priority", MEDIUM, "priority of the task LOW|MEDIUM|HIGH")
	status := flag.String("status", TODO, "status of the tasks TODO|INPROGRESS|DONE")
	flag.Var(&startDate, "start", "start date")
	flag.Var(&endDate, "end", "end date")
	flag.Parse()

	switch *action {
	case LIST, CREATE, UPDATE, DELETE:
	default:
		flag.Usage()
		os.Exit(2)
	}
	switch *priority {
	case LOW, MEDIUM, HIGH:
	default:
		flag.Usage()
		os.Exit(2)
	}
	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		log.Fatal()
	}
	defer conn.Close()

	client := api.NewToDoApiClient(conn)

	switch *action {
	case LIST:
		tasks, err := client.List(context.Background(), &api.Empty{})
		if err != nil {
			log.Fatal(err)
		}
		for _, task := range tasks.Tasks {
			fmt.Println(task)
		}
	case CREATE:
		t := &api.Task{}
		t.Uid = int32(*uid)
		t.Start = startDate.Unix()
		t.End = endDate.Unix()
		t.Subject = *subject
		t.Status = api.Task_Status(api.Task_Status_value[*status])
		t.Priority = api.Task_Priority(api.Task_Priority_value[*priority])
		ret, err := client.Create(context.Background(), t)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(ret.GetType())
	default:
	}
}

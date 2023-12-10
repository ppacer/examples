package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/dskrzypiec/scheduler/dag"
	"github.com/dskrzypiec/scheduler/db"
	"github.com/dskrzypiec/scheduler/exec"
	"github.com/dskrzypiec/scheduler/meta"
	"github.com/dskrzypiec/scheduler/sched"
)

type EmptyTask struct {
	TaskId string
}

func (et EmptyTask) Id() string { return et.TaskId }
func (et EmptyTask) Execute() {
	fmt.Printf(" ========== EmptyTask: %s ==========\n", et.TaskId)
}

func setupDAGs() {
	//        t2 ------
	//      /          \
	// root              end
	//      \          /
	//        t3 --> t4
	root := dag.Node{Task: EmptyTask{TaskId: "Start"}}
	t2 := dag.Node{Task: EmptyTask{TaskId: "t2"}}
	t3 := dag.Node{Task: EmptyTask{TaskId: "t3"}}
	t4 := dag.Node{Task: EmptyTask{TaskId: "t4"}}
	end := dag.Node{Task: EmptyTask{TaskId: "Finish"}}
	root.Next(&t2)
	root.Next(&t3)
	t3.Next(&t4)
	t2.Next(&end)
	t4.Next(&end)

	startTs := time.Date(2023, time.December, 8, 12, 0, 0, 0, time.UTC)
	schedule := dag.FixedSchedule{Interval: 5 * time.Second, Start: startTs}
	emptyDag := dag.New(dag.Id("empty_dag")).
		AddSchedule(&schedule).
		AddRoot(&root).
		Done()

	// Add emptyDag to the central DAG repository.
	dag.Add(emptyDag)
}

//go:embed *.go
var taskGoFiles embed.FS

func main() {
	const port = 8080
	setupDAGs()
	meta.ParseASTs(taskGoFiles)

	dbClient, dbErr := db.NewSqliteClient("scheduler.db")
	if dbErr != nil {
		log.Panic(dbErr)
	}
	scheduler := sched.New(dbClient)
	schedulerHttpHandler := scheduler.Start()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: schedulerHttpHandler,
	}

	// Run executor within the same program
	go func() {
		executor := exec.New(fmt.Sprintf("http://localhost:%d", port), nil)
		executor.Start()
	}()

	lasErr := server.ListenAndServe()
	if lasErr != nil {
		slog.Error("ListenAndServer failed", "err", lasErr)
		log.Panic("Cannot start the server")
	}
}

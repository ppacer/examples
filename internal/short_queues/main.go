package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/ppacer/core/dag"
	"github.com/ppacer/core/db"
	"github.com/ppacer/core/exec"
	"github.com/ppacer/core/meta"
	"github.com/ppacer/core/scheduler"
)

type EmptyTask struct {
	TaskId string
}

func (et EmptyTask) Id() string { return et.TaskId }
func (et EmptyTask) Execute() {
	fmt.Printf(" ========== EmptyTask: %s ==========\n", et.TaskId)
}

func emptyLinkedList(dagId string, len int, interval time.Duration) dag.Dag {
	// Start -> step_1 -> step_2 -> ... step_{len - 1}
	root := dag.Node{Task: EmptyTask{TaskId: "Start"}}
	prev := &root
	for i := 1; i < len-1; i++ {
		n := dag.Node{Task: EmptyTask{TaskId: fmt.Sprintf("step_%d", i)}}
		prev.Next(&n)
		prev = &n
	}

	startTs := time.Date(2024, time.January, 14, 8, 0, 0, 0, time.UTC)
	schedule := dag.FixedSchedule{Interval: interval, Start: startTs}
	linkedList := dag.New(dag.Id(dagId)).
		AddSchedule(&schedule).
		AddRoot(&root).
		Done()
	return linkedList
}

//go:embed *.go
var taskGoFiles embed.FS

func main() {
	const port = 8080
	// Add emptyDag to the central DAG repository.
	dag.Add(emptyLinkedList("short_ll", 10, 30*time.Second))
	dag.Add(emptyLinkedList("short_ll_freq", 10, 15*time.Second))
	dag.Add(emptyLinkedList("longer_ll", 100, 30*time.Second))

	meta.ParseASTs(taskGoFiles)

	dbClient, dbErr := db.NewSqliteClient("scheduler.db")
	if dbErr != nil {
		log.Panic(dbErr)
	}
	config := scheduler.DefaultConfig
	config.DagRunQueueLen = 1
	config.DagRunTaskQueueLen = 1
	fmt.Printf("config: %+v\n", config)
	scheduler := scheduler.New(dbClient, scheduler.DefaultQueues(config), config)
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

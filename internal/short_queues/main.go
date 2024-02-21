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
func (et EmptyTask) Execute(tc dag.TaskContext) {
	fmt.Printf(" ========== EmptyTask: %s ==========\n", et.TaskId)
	tc.Logger.Warn("Empty task finished successfully!", "ts", time.Now())
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

func prepDagRegistry() dag.Registry {
	shortLL := emptyLinkedList("short_ll", 10, 30*time.Second)
	shortLLFreq := emptyLinkedList("short_ll_freq", 10, 15*time.Second)
	longerLL := emptyLinkedList("longer_ll", 100, 30*time.Second)

	return dag.Registry{
		shortLL.Id:     shortLL,
		shortLLFreq.Id: shortLLFreq,
		longerLL.Id:    longerLL,
	}
}

func setupSQLiteDBs() (*db.Client, *db.Client) {
	dbClient, dbErr := db.NewSqliteClient("scheduler.db")
	if dbErr != nil {
		log.Panic(dbErr)
	}
	logsDbClient, logsDbErr := db.NewSqliteClientForLogs("logs.db")
	if logsDbErr != nil {
		log.Panic(logsDbErr)
	}
	return dbClient, logsDbClient
}

//go:embed *.go
var taskGoFiles embed.FS

func main() {
	const port = 8080
	dags := prepDagRegistry()

	meta.ParseASTs(taskGoFiles)

	dbClient, logsDbClient := setupSQLiteDBs()
	config := scheduler.DefaultConfig
	config.DagRunQueueLen = 1
	config.DagRunTaskQueueLen = 1
	fmt.Printf("config: %+v\n", config)
	scheduler := scheduler.New(dbClient, scheduler.DefaultQueues(config), config)
	schedulerHttpHandler := scheduler.Start(dags)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: schedulerHttpHandler,
	}

	// Run executor within the same program
	go func() {
		executor := exec.New(fmt.Sprintf("http://localhost:%d", port), logsDbClient, nil)
		executor.Start(dags)
	}()

	lasErr := server.ListenAndServe()
	if lasErr != nil {
		slog.Error("ListenAndServer failed", "err", lasErr)
		log.Panic("Cannot start the server")
	}
}

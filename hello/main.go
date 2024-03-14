package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/ppacer/core/dag"
	"github.com/ppacer/core/db"
	"github.com/ppacer/core/exec"
	"github.com/ppacer/core/meta"
	"github.com/ppacer/core/scheduler"
)

//go:embed *.go
var taskGoFiles embed.FS

func main() {
	const port = 9321
	meta.ParseASTs(taskGoFiles)
	printDag := printDAG("example")
	dags := dag.Registry{}
	dags[printDag.Id] = printDag

	// Setup default scheduler
	schedulerServer := scheduler.DefaultStarted(dags, "scheduler.db", port)

	// Setup and run executor in a separate goroutine
	go func() {
		schedUrl := fmt.Sprintf("http://localhost:%d", port)
		logsDbClient, logsDbErr := db.NewSqliteClientForLogs("logs.db", nil)
		if logsDbErr != nil {
			log.Panic(logsDbErr)
		}
		executor := exec.New(schedUrl, logsDbClient, nil, nil)
		executor.Start(dags)
	}()

	// Start scheduler HTTP server
	lasErr := schedulerServer.ListenAndServe()
	if lasErr != nil {
		slog.Error("ListenAndServer failed", "err", lasErr)
		log.Panic("Cannot start the server")
	}
}

type printTask struct {
	taskId string
}

func (pt printTask) Id() string { return pt.taskId }

func (pt printTask) Execute(tc dag.TaskContext) error {
	fmt.Printf(" >>> PrintTask <<<: %s\n", pt.taskId)
	tc.Logger.Info("PrintTask finished!", "ts", time.Now())
	return nil
}

func printDAG(dagId string) dag.Dag {
	// [start] --> [end]
	start := dag.NewNode(printTask{taskId: "start"})
	start.NextTask(printTask{taskId: "finish"})

	startTs := time.Date(2024, time.March, 11, 12, 0, 0, 0, time.UTC)
	schedule := dag.FixedSchedule{Interval: 10 * time.Second, Start: startTs}

	printDag := dag.New(dag.Id(dagId)).
		AddSchedule(&schedule).
		AddRoot(start).
		Done()
	return printDag
}

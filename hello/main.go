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

func emptyDag() dag.Dag {
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
	schedule := dag.FixedSchedule{Interval: 15 * time.Second, Start: startTs}
	emptyDag := dag.New(dag.Id("empty_dag")).
		AddSchedule(&schedule).
		AddRoot(&root).
		Done()
	return emptyDag
}

//go:embed *.go
var taskGoFiles embed.FS

func main() {
	const port = 8080
	eDag := emptyDag()
	dags := dag.Registry{
		eDag.Id: eDag,
	}

	meta.ParseASTs(taskGoFiles)

	dbClient, dbErr := db.NewSqliteClient("scheduler.db")
	if dbErr != nil {
		log.Panic(dbErr)
	}
	logsDbClient, logsDbErr := db.NewSqliteClientForLogs("logs.db")
	if logsDbErr != nil {
		log.Panic(logsDbErr)
	}
	config := scheduler.DefaultConfig
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

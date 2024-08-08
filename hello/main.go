package main

import (
	"context"
	"time"

	"github.com/ppacer/core"
	"github.com/ppacer/core/dag"
	"github.com/ppacer/core/dag/schedule"
	"github.com/ppacer/tasks"
	"github.com/ppacer/ui"
)

const (
	schedulerPort = 9321
	uiPort        = 9322
)

func main() {
	ctx := context.Background()

	dags := dag.Registry{}
	dags.Add(printDAG("hello_world_dag"))

	go func() {
		ui.DefaultStarted(schedulerPort, uiPort)
	}()
	core.DefaultStarted(ctx, dags, schedulerPort)
}

func printDAG(dagId string) dag.Dag {
	//         t21
	//       /
	// start
	//       \
	//         t22 --> finish
	start := dag.NewNode(tasks.NewPrintTask("start", "hello"))
	t21 := dag.NewNode(tasks.NewPrintTask("t21", "foo"))
	t22 := dag.NewNode(tasks.NewPrintTask("t22", "bar"))
	finish := dag.NewNode(tasks.NewPrintTask("finish", "I'm done!"))

	start.Next(t21)
	start.Next(t22)
	t22.Next(finish)

	startTs := time.Date(2024, time.March, 11, 12, 0, 0, 0, time.Local)
	schedule := schedule.NewFixed(startTs, 10*time.Second)

	printDag := dag.New(dag.Id(dagId)).
		AddSchedule(&schedule).
		AddRoot(start).
		Done()
	return printDag
}

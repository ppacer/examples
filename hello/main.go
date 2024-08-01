package main

import (
	"context"
	"time"

	"github.com/ppacer/core"
	"github.com/ppacer/core/dag"
	"github.com/ppacer/core/dag/schedule"
	"github.com/ppacer/tasks"
)

func main() {
	const port = 9321
	ctx := context.Background()
	dags := dag.Registry{}
	dags.Add(printDAG("printing_dag"))
	core.DefaultStarted(ctx, dags, port)
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

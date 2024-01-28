package twoBinaries

import (
	"embed"
	"fmt"
	"time"

	"github.com/ppacer/core/dag"
)

//go:embed *.go
var TaskGoFiles embed.FS

type EmptyTask struct {
	DagId  string
	TaskId string
	N      int
	Delay  time.Duration
}

func (et EmptyTask) Id() string { return et.TaskId }
func (et EmptyTask) Execute() {
	fmt.Printf(" ========== Starting: %s ==========\n", et.TaskId)
	for i := 0; i < et.N; i++ {
		fmt.Printf("Calculating f(%s, %s, %d)...\n", et.DagId, et.TaskId, i)
		time.Sleep(et.Delay)
	}
	fmt.Printf(" ========== Finished: %s ==========\n", et.TaskId)
}

func emptyLinkedList(dagId string, len int, interval time.Duration) dag.Dag {
	// Start -> step_1 -> step_2 -> ... step_{len - 1}
	w := 3 * time.Second
	root := dag.Node{Task: EmptyTask{DagId: dagId, TaskId: "Start", N: 10, Delay: w}}
	prev := &root
	for i := 1; i < len-1; i++ {
		n := dag.Node{Task: EmptyTask{
			DagId:  dagId,
			TaskId: fmt.Sprintf("step_%d", i),
			N:      5,
			Delay:  w,
		}}
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

func SetupDags() dag.Registry {
	d1 := emptyLinkedList("dag_1", 25, 5*time.Minute)
	d2 := emptyLinkedList("dag_2", 5, 1*time.Minute+15*time.Second)
	return dag.Registry{
		d1.Id: d1,
		d2.Id: d2,
	}
}

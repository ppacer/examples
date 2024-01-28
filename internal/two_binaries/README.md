# Scheduler and Executor in two binaries

The only requirements to run this example program is having Go in version
`>=1.21` and SQLite. To compile and run both programs, just do the following:


```
go build -o sched ./cmd/sched
go build -o exec ./cmd/exec
trap 'kill %1; kill %2' SIGINT; ./sched & ./exec & wait
```

It builds two binaries, one for the scheduler and the other one for the
executor. Command `trap` makes it easy to kill both programs when you click
`Ctrl+C`.


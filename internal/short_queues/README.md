# Short internal queues

This example contains ppacer setup in which `DagRunQueue` and `DagRunTaskQueue`
lengths are set to one. This is to show the corner case when there is basically
no parallelisation and its performance.

```
go generate
go build
./hello
```

Currently there is only backend side - scheduler, executor and database. If you
want to checkout details on your DAG runs, please explore `scheduler.db` SQLite
database. Few examples:

```
sqlite3 scheduler.db 'SELECT DagId, StartTs, Schedule, CreateTs FROM dags'
sqlite3 scheduler.db 'SELECT DagId, TaskId, IsCurrent, InsertTs, TaskTypeName FROM dagtasks'
sqlite3 scheduler.db 'SELECT DagId, TaskId, IsCurrent, TaskBodySource FROM dagtasks'
sqlite3 scheduler.db 'SELECT * FROM dagruns'
sqlite3 scheduler.db 'SELECT * FROM dagruntasks'
```

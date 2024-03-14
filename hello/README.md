# Hello World example

The only requirements to run hello example program is having Go in version
`>=1.22` and SQLite. To compile and run the example just run:

```
go generate
go build
./ppacer_demo
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

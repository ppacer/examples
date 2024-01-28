package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"twoBinaries"

	"github.com/ppacer/core/db"
	"github.com/ppacer/core/meta"
	"github.com/ppacer/core/scheduler"
)

func main() {
	const port = 8080
	meta.ParseASTs(twoBinaries.TaskGoFiles)
	dags := twoBinaries.SetupDags()

	dbClient, dbErr := db.NewSqliteClient("scheduler.db")
	if dbErr != nil {
		log.Panic(dbErr)
	}
	config := scheduler.DefaultConfig
	fmt.Printf("config: %+v\n", config)
	scheduler := scheduler.New(dbClient, scheduler.DefaultQueues(config), config)
	schedulerHttpHandler := scheduler.Start(dags)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: schedulerHttpHandler,
	}

	lasErr := server.ListenAndServe()
	if lasErr != nil {
		slog.Error("ListenAndServer failed", "err", lasErr)
		log.Panic("Cannot start the server")
	}
}

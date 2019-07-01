package monitor

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/parser"
)

func monitorQueries() {
	queries, err := db.GetQueries()
	if err != nil {
		panic("error")
	}

	for _, query := range *queries {
		gocron.Every(uint64(query.MonitoringInterval)).Minutes().Do(parser.Run)
	}
	<-gocron.Start()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		gocron.Clear()
		os.Exit(1)
	}()
}

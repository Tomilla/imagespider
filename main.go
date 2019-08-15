package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Tomilla/imagespider/t66y/generator"

	"github.com/Tomilla/imagespider/config"

	_ "github.com/Tomilla/imagespider/all"
	"github.com/Tomilla/imagespider/scheduler"

	"github.com/Tomilla/imagespider/engine"
)

func main() {

	e := engine.NewConcurrentEngine(config.C.GetImageChan())

	e.Run(scheduler.NewScheduler(), generator.NewGenerator(config.C.GetStartPages()))

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
	{
		e.Shutdown()
	}
}

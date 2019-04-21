package main

import (
	"github.com/wuxiangzhou2010/imagespider/t66y/generator"
	"os"
	"os/signal"
	"syscall"

	"github.com/wuxiangzhou2010/imagespider/config"

	_ "github.com/wuxiangzhou2010/imagespider/all"
	"github.com/wuxiangzhou2010/imagespider/scheduler"

	"github.com/wuxiangzhou2010/imagespider/engine"
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

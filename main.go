package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/wuxiangzhou2010/imagespider/config"

	_ "github.com/wuxiangzhou2010/imagespider/all"
	"github.com/wuxiangzhou2010/imagespider/scheduler"

	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler/util/agent/my"
	"github.com/wuxiangzhou2010/imagespider/engine"
	"github.com/wuxiangzhou2010/imagespider/t66y/parser"
)

func main() {

	e := engine.NewConcurrentEngine(config.C.GetImageChan())
	e.Run(scheduler.NewScheduler(), generateStartPages())

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
	{
		e.Shutdown()
	}
}

func generateStartPages() (r []engine.Request) {

	for _, url := range config.C.GetStartPages() {
		r = append(r, engine.Request{
			Url:        url,
			ParserFunc: parser.ParseTopicList,
			Agent:      my.NewAgent(),
		})
	}
	return

}

func init() {
	go http.ListenAndServe(":6060", nil)
}

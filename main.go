package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/all"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/scheduler"

	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler/util/agent/my"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/engine"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/t66y/parser"
)

func main() {
	engine.Run(scheduler.NewScheduler(), generatedStartRequest()...)

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}

func generatedStartRequest() (r []engine.Request) {

	for _, url := range config.StartPages {
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

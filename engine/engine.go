package engine

import (
	"time"

	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/fetcher"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"

	"log"
)

type ConcurrentEngine struct {
	ImageChan chan model.Topic
	s         Scheduler
}

func NewConcurrentEngine(imageChan chan model.Topic) *ConcurrentEngine {
	return &ConcurrentEngine{ImageChan: imageChan}
}

func (e *ConcurrentEngine) Shutdown() {
	close(e.ImageChan) // 不继续发送图片下载
	e.s.Shutdown()
	time.Sleep(10)

}
func (e *ConcurrentEngine) Run(s Scheduler, seeds []Request) {
	out := make(chan ParseResult)
	e.s = s

	go s.Schedule() // scheduler started
	w := newWorker()
	for i := 0; i < s.GetWorkCount(); i++ {
		go w.work(s, out) // 创建所有worker
	}

	for _, r := range seeds { // submit start page
		s.SubmitRequest(r)
	}

	for {
		result := <-out
		for _, r := range result.Requests {
			s.SubmitRequest(r)
		}

		e.dealItems(result.Items)
	}

}

type Worker struct {
}

func newWorker() *Worker {
	return &Worker{}
}

// fetch as request and return the parsed result

func (w *Worker) work(s Scheduler, out chan ParseResult) {
	workChan := make(chan Request)
	s.SubmitWorker(workChan)

	for {
		r := <-workChan

		log.Printf("Fetching %s\n", r.Url)
		body, err := fetcher.Fetch(r.Url)
		if err != nil {
			continue
		}

		ParseResult := r.ParserFunc(body)
		out <- ParseResult
		s.SubmitWorker(workChan)
	}

}

// deal all items that need not fetch again
func (e *ConcurrentEngine) dealItems(items []interface{}) {
	for _, item := range items {

		switch item.(type) {
		case model.Topic:
			if imageChan := e.ImageChan; imageChan != nil {
				imageChan <- item.(model.Topic) // 转换为topic 类型
			}
		default:
			log.Printf("Got item %s", item)
		}
	}
}

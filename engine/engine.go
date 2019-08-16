package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/Tomilla/imagespider/config"
	"github.com/Tomilla/imagespider/model"
)

type ConcurrentEngine struct {
	ImageChan   chan model.Topic
	s           Scheduler
	ElasticChan chan model.Topic
}

func NewConcurrentEngine(imageChan chan model.Topic) *ConcurrentEngine {
	return &ConcurrentEngine{ImageChan: imageChan}
}

func (e *ConcurrentEngine) Shutdown() {
	close(e.ImageChan) // 不继续发送图片下载
	e.s.Shutdown()
	// e.ElasticChan.Stop() // stop elasticSearch client
	time.Sleep(10)

}
func (e *ConcurrentEngine) Run(s Scheduler, requestChan chan Request) {

	e.s = s
	e.ElasticChan = config.C.GetElasticChan() // new ElasticChan client

	out := make(chan ParseResult)
	hungry := make(chan bool)
	go s.Schedule(hungry) // scheduler started

	w := newWorker()
	for i := 0; i < config.C.GetEngineWorkerCount(); i++ {
		go w.work(s, out) // 创建所有worker
	}

	go func() {
		for {
			select {
			case <-hungry: // 请求下一页
				r, more := <-requestChan
				if more {
					fmt.Println("Got next page, ", r.Url)
					s.SubmitRequest(r)
				} else {
					fmt.Println("All initial pages are sent")
					return
				}
			case <-e.ElasticChan:
				r := <-e.ElasticChan
				fmt.Println("ES:", r)
			}
		}
	}()

	for {
		select {
		case result := <-out: // 页面解析结果
			for _, r := range result.Requests {
				go s.SubmitRequest(r)
			}

			e.dealItems(result.Items)
		}
	}

}

// deal all items that need not fetch again
func (e *ConcurrentEngine) dealItems(items []interface{}) {
	for _, item := range items {

		switch item.(type) {
		case model.Topic:
			if imageChan := e.ImageChan; imageChan != nil {
				imageChan <- item.(model.Topic) // save to image
			}
			if e.ElasticChan != nil {
				e.ElasticChan <- item.(model.Topic) // save to elasticSearch
			}

		default:
			log.Printf("[engine dealItems] Got item %s", item)
		}
	}
}

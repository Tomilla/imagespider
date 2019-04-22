package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/wuxiangzhou2010/imagespider/config"
	"github.com/wuxiangzhou2010/imagespider/model"
	"gopkg.in/olivere/elastic.v5"
)

type ConcurrentEngine struct {
	ImageChan chan model.Topic
	s         Scheduler
	elastic   *elastic.Client
}

func NewConcurrentEngine(imageChan chan model.Topic) *ConcurrentEngine {
	return &ConcurrentEngine{ImageChan: imageChan}
}

func (e *ConcurrentEngine) Shutdown() {
	close(e.ImageChan) // 不继续发送图片下载
	e.s.Shutdown()
	e.elastic.Stop() // stop elasticSearch client
	time.Sleep(10)

}
func (e *ConcurrentEngine) Run(s Scheduler, requestChan chan Request) {

	e.s = s
	e.elastic = NewConnection() // new elastic client

	out := make(chan ParseResult)
	hungry := make(chan bool)
	go s.Schedule(hungry) // scheduler started

	w := newWorker()
	for i := 0; i < config.C.GetEngineWorkerCount(); i++ {
		go w.work(s, out) // 创建所有worker
	}

	for {
		select {
		case <-hungry: // 请求下一页
			r := <-requestChan
			fmt.Println("Got next page, ", r.Url)
			s.SubmitRequest(r)
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
				imageChan <- item.(model.Topic) // 转换为topic 类型
			}

			if e.elastic != nil {
				e.saveElasticSearch(item.(model.Topic))
			}

		default:
			log.Printf("[engine dealItems ]Got item %s", item)
		}
	}
}

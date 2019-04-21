package engine

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/wuxiangzhou2010/imagespider/config"
	"time"

	"github.com/wuxiangzhou2010/imagespider/fetcher"
	"github.com/wuxiangzhou2010/imagespider/model"
	"gopkg.in/olivere/elastic.v5"
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
func (e *ConcurrentEngine) Run(s Scheduler, requestChan chan Request) {

	e.s = s

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
				s.SubmitRequest(r)
			}

			e.dealItems(result.Items)
		}
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

		log.Printf("Fetching %s, %s \n", r.Name, r.Url)
		body, err := fetcher.Fetch(r.Url)
		if err != nil {
			//panic(err)
			log.Println("Fetching error:", err, r.Url)
			s.SubmitWorker(workChan)
			continue
		}

		ParseResult := r.ParserFunc(body, r.Url)
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
			} else {
				e.saveElasticSearch(item.(model.Topic))
			}
		default:
			log.Printf("Got item %s", item)
		}
	}
}

func (e *ConcurrentEngine) saveElasticSearch(topic model.Topic) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	resp, err := client.Index().
		Index("t66y").
		Type("topics").Id(hash(topic.Url)).
		BodyJson(topic).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf(" %+v\n", resp)

}

func hash(s string) string {

	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	fmt.Println(s)
	return hex.EncodeToString(bs)

}

package engine

import (
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/downloader"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/fetcher"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"
	"log"
)

func Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}
	for len(requests) > 0 {
		r := requests[0]

		result, err := work(&r)
		if err != nil {
			continue
		}
		requests = requests[1:] // 成功后删除第一个request

		requests = append(requests, result.Requests...)
		dealItems(result.Items)

	}
}

func work(r *Request) (*ParseResult, error) {

	log.Printf("Fetching %s\n", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher : error "+"fetching url %s: %v", r.Url, err)
		return nil, err
	}

	ParseResult := r.ParserFunc(body)
	return &ParseResult, nil
}

// deal all items
func dealItems(items []interface{}) {
	for _, item := range items {

		switch item.(type) {
		case model.Topic:
			downloader.ImageChan <- item.(model.Topic)
		default:
			log.Printf("Got item %s", item)
		}
	}
}

package main

import (
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/downloader"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/engine"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/t66y/parser"
)

func main() {
	engine.Run(engine.Request{
		Url:        "http://t66y.com/thread0806.php?fid=8",
		ParserFunc: parser.ParseTopicList,
	})
}

func init() {
	downloader.NewWorker(parser.ImageCh)
}

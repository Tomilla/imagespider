package downloader

import (
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"
)

var count int32

type Worker struct {
	*config.Config
	workChan chan model.Topic
}

type Downloader  interface{
	Download(chan interface{})
}


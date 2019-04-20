package image

import (
	"fmt"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"
)

var count int32

func init() {

	ch := make(chan model.Topic)

	config.C.SetImageChan(ch) // engine 通过 这个channel 和downloader 通信

	d := NewDownloader(config.C.GetImageConfig())

	go d.Run() // start downloader

	fmt.Println("Image init")

}

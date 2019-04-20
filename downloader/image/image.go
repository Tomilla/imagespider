package image

import (
	"fmt"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"
)

var count int32

func init() {

	ch := make(chan model.Topic)
	config.C.SetImageChan(ch)
	d := NewDownloader(config.C.GetImageConfig())
	go d.Run()

	fmt.Println("Image init")

}

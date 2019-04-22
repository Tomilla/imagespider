package image

import (
	"github.com/wuxiangzhou2010/imagespider/config"
	"github.com/wuxiangzhou2010/imagespider/model"
)

var count int32

func init() {

	ch := make(chan model.Topic)
	hungryChan :=make(chan bool)
	config.C.SetImageChan(ch) // engine 通过 这个channel 和downloader 通信
	config.C.SetImageHungryChan(hungryChan) // engine 通过 这个channel 和downloader 通信

	d := NewDownloader(config.C.GetImageConfig())

	go d.Run() // start downloader

	//fmt.Println("Image init")

}

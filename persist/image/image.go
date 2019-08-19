package image

import (
    "github.com/Tomilla/imagespider/common"
    "github.com/Tomilla/imagespider/model"
)

var count int32

func init() {

    ch := make(chan model.Topic)
    hungryChan := make(chan bool)
    common.C.SetImageChan(ch)               // engine 通过 这个channel 和downloader 通信
    common.C.SetImageHungryChan(hungryChan) // engine 通过 这个channel 和downloader 通信

    d := NewDownloader(common.C.GetImageConfig(), common.C.GetLimitConfig())

    go d.Run() // start persist

    // fmt.Println("Image init")

}

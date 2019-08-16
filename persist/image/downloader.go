package image

import (
    "os"
    "path"
    "strconv"

    "github.com/Tomilla/imagespider/config"
    "github.com/Tomilla/imagespider/util"
)

type downloader struct {
    config.ImageConfig
    workChan chan work
}

func NewDownloader(imageConfig *config.ImageConfig) *downloader {
    return &downloader{ImageConfig: *imageConfig}
}

func (d *downloader) Run() {
    if err := os.MkdirAll(d.Path, util.DefaultFilePerm); err != nil {
        panic(err)
    }
    readyChan := make(chan chan work)
    workChan := make(chan work)
    s := newScheduler(workChan, readyChan)
    go s.schedule()

    d.CreateWorker(s)
    for {
        topic := <-d.ImageChan

        // 分析 url 和名字
        baseFolder := path.Join(d.Path, topic.Name)
        // fmt.Println("BaseFolder", baseFolder)
        if !d.UniqFolder { // 如果不是统一文件夹， 则需要分别创建文件夹
            if err := os.MkdirAll(baseFolder, util.DefaultFilePerm); err != nil {
                panic(err)
            }
        }

        for i, url := range topic.Images { // pass to worker
            fileName := d.getFileName(baseFolder, topic.Name, i)
            w := newWork(url, fileName)
            workChan <- w
        }

    }
}

func (d *downloader) CreateWorker(s *scheduler) {
    ws := NewWorkers(s)
    ws.Start()
}

func (d *downloader) getFileName(baseFolder, name string, index int) string {
    if d.UniqFolder { // only provide filename is fine for unique folder
        return baseFolder + strconv.Itoa(index) + ".jpg"
    }
    // otherwise, we need join the post name with filename
    return path.Join(baseFolder, name+strconv.Itoa(index)+".jpg")

}

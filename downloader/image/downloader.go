package image

import (
	"fmt"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"os"
	"path"
	"strconv"
)

type downloader struct {
	config.ImageConfig
	workChan chan work
}

func NewDownloader(imageConfig *config.ImageConfig) *downloader {

	return &downloader{ImageConfig: *imageConfig}
}
func (d *downloader) Run() {

	if err := os.MkdirAll(d.Path, 0700); err != nil {
		panic(err)
	}
	readyChan := make(chan chan work)
	workChan := make(chan work)
	s := newScheduler(workChan, readyChan)
	go s.schedule()

	d.CreateWorker(s)
	for {
		topic := <-d.ImageChan

		baseFolder := path.Join(d.Path, topic.Name)
		fmt.Println("BaseFolder", baseFolder)
		if !d.UniqFolder { // 如果不是统一文件夹， 则需要分别创建文件夹
			if err := os.MkdirAll(baseFolder, 0700); err != nil {
				panic(err)
			}
		}

		for i, url := range topic.Images {
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

// golang新版本的应该
func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (d *downloader) getFileName(baseFolder, name string, index int) string {
	if d.UniqFolder { //放同一个文件夹， 只需要提供文件名就行
		return baseFolder + strconv.Itoa(index) + ".jpg"
	}
	return path.Join(baseFolder, name+strconv.Itoa(index)+".jpg") // 放

}

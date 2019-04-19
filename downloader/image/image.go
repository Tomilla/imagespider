package image

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"
)

var count int32

type Worker struct {
	Downloader
}

type Downloader struct {
	config.ImageConfig
}

func NewDownloader(imageConfig config.ImageConfig) *Downloader {
	return &Downloader{ImageConfig: imageConfig}
}

func (d *Downloader) Download() {

	if err := os.MkdirAll(d.Path, 0700); err != nil {
		panic(err)
	}
	d.CreateWorker()
}

func (d *Downloader) CreateWorker() {
	w := &Worker{
		*d,
	}

	for i := 0; i < d.WorkerCount; i++ {

		go w.Work()
	}
	fmt.Printf("create %d image worker\n", d.WorkerCount)
}

func (w *Worker) Work() {
	for {
		tp, ok := <-w.ImageChan
		if !ok {
			return // channel 关闭，退出
		}

		w.Download(tp)

	}

}

func (w *Worker) Download(tp model.Topic) {
	baseFolder := path.Join(w.Path, tp.Name)
	fmt.Println("Basefolder", baseFolder)
	if !w.UniqFolder { // 如果不是统一文件夹， 则需要分别创建文件夹
		if err := os.MkdirAll(baseFolder, 0700); err != nil {
			panic(err)
		}
	}

	for i, url := range tp.Images {
		err := w.downloadWithPath(url, baseFolder, tp.Name, i)
		if err != nil {
			log.Println("####### Error download ", err, url)
			fineName := w.getFileName(baseFolder, tp.Name, i)
			os.Remove(fineName) // 下载失败 删除文件

			continue
		}
		//log.Printf("#%d downloaded %s", atomic.AddInt32(&count, int32(len(tp.Images))), tp.Name)
		log.Printf("#%d downloaded %s", atomic.AddInt32(&count, 1), tp.Name)

	}

}
func (w *Worker) downloadWithPath(link, baseFolder, name string, index int) error {
	fileName := w.getFileName(baseFolder, name, index)
	if pathExist(fileName) {
		return nil
	}
	//resp, err := http.Get(link)
	//@@@@@@@@@@@@@@@@@

	proxyStr := "socks5://localhost:1080"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}
	tr := &http.Transport{ //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr, //解决x509: certificate signed by unknown authority
	}
	req, err := http.NewRequest("GET", link, nil)
	resp, err := client.Do(req)

	//@@@@@@@@@@@@@@@@@

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf := bufio.NewReader(resp.Body)
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}

	io.Copy(out, buf)
	defer out.Close()
	return nil
}

func (w *Worker) getFileName(baseFolder, name string, index int) string {
	if w.UniqFolder {
		return baseFolder + strconv.Itoa(index) + ".jpg"
	}
	return baseFolder + "/" + name + strconv.Itoa(index) + ".jpg"

}

// golang新版本的应该
func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {

	imageConf := config.NewImageConfig()
	imageConf.ImageChan = make(chan model.Topic)
	d := NewDownloader(*imageConf)
	go d.Download()
	fmt.Println("Image init")
	config.AllConf.Image = *imageConf
}

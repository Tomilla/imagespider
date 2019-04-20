package image

import (
	"bufio"
	"crypto/tls"

	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"sync/atomic"
	"time"
)

type Worker struct {
	workChan    chan work
	workerCount int
}

func NewWorker(workChan chan work, workerCount int) *Worker {
	return &Worker{workChan: workChan,
		workerCount: workerCount}
}

type work struct {
	url      string
	fileName string
}

func newWork(url string, fileName string) work {
	return work{url: url, fileName: fileName}
}

func (w *Worker) Start() {

	for i := 0; i < w.workerCount; i++ {
		go w.work()
	}

}
func (w *Worker) work() {
	for {
		task, ok := <-w.workChan
		if !ok {
			return // channel 关闭，退出
		}

		w.Download(task)

	}

}

func (w *Worker) Download(task work) {

	err := w.downloadWithPath(task.url, task.fileName)
	if err != nil {
		log.Println("####### Error download ", err, task.url)
		os.Remove(task.fileName) // 下载失败 删除文件
	}
	//log.Printf("#%d downloaded %s", atomic.AddInt32(&count, int32(len(tp.Images))), tp.Name)
	log.Printf("#%d downloaded %s", atomic.AddInt32(&count, 1), task.fileName)

}

func (w *Worker) downloadWithPath(link, fileName string) error {

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

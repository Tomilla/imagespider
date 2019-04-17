package downloader

import (
	"crypto/tls"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Worker struct {
	*config.Config
	workCh chan []*model.Topic
}

func NewWorker(workCh chan []*model.Topic) {
	w := &Worker{
		config.DefaultConfig,
		workCh,
	}

	for i := 0; i < 20; i++ {
		go w.Work()
	}
}

func (w *Worker) Work() {
	for {
		tps := <-w.workCh
		for _, t := range tps {
			w.download(t)
		}
	}

}

func (w *Worker) download(tp *model.Topic) {
	baseFolder := w.Config.Path + "/" + tp.Name
	if err := os.MkdirAll(baseFolder, 0700); err != nil {
		panic(err)
	}
	for i, url := range tp.Images {
		err := w.downloadWithPath(url, baseFolder, tp.Name, i)
		if err != nil {

			log.Println("#######Error download ", err, url)
			//panic(err)
			continue
		}
	}
	log.Println("downloaded ", tp.Name, "-->", len(tp.Images))
}
func (w *Worker) downloadWithPath(url, baseFolder, name string, index int) error {
	//log.Println("Download link ", url)
	fileName := w.getFileName(baseFolder, name, index)
	//fmt.Println(fileName)
	if pathExist(fileName) {
		return nil
	}
	//resp, err := http.Get(url)
	//@@@@@@@@@@@@@@@@@
	tr := &http.Transport{ //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr, //解决x509: certificate signed by unknown authority
	}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	//@@@@@@@@@@@@@@@@@

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}

	io.Copy(out, resp.Body)
	defer out.Close()
	return nil
}

func (w *Worker) getFileName(baseFolder, name string, index int) string {
	if w.Config.SeperateFolder {
		return baseFolder + "/" + name + strconv.Itoa(index) + ".jpg"
	}
	return baseFolder + strconv.Itoa(index) + ".jpg"

}

// golang新版本的应该
func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

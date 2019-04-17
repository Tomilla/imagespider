package downloader

import (
	"fmt"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"
	"io"
	"net/http"
	"os"
	"strconv"
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

	for i := 0; i < 10; i++ {
		go w.Work()
	}
}

func (w *Worker) Work() {
	tps := <-w.workCh
	for _, t := range tps {
		w.download(t)
	}

}

func (w *Worker) download(tp *model.Topic) {
	baseFolder := w.Config.Path + "/" + tp.Name
	if err := os.MkdirAll(baseFolder, 0700); err != nil {
		panic(err)
	}
	fmt.Println("Saving image file : ", baseFolder)
	for i, url := range tp.Images {
		err := w.downloadWithPath(url, baseFolder, tp.Name, i)
		if err != nil {
			fmt.Println("Error download ", url)
			continue
		}
	}
}
func (w *Worker) downloadWithPath(url, baseFolder, name string, index int) error {
	fmt.Println("Download link ", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fileName := w.getFileName(baseFolder, name, index)
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

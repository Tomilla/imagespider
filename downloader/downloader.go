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
}

func NewWorker() *Worker {
	return &Worker{
		config.DefaultConfig,
	}
}

func (w *Worker) Work(tps []*model.Topic) {

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
	for i, l := range tp.Iamges {
		fmt.Println("Download link ", l)
		resp, err := http.Get(l)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		file := baseFolder + "/" + strconv.Itoa(i) + ".jpg"

		out, err := os.Create(file)
		if err != nil {
			panic(err)
		}

		io.Copy(out, resp.Body)
		defer out.Close()
	}
}

package image

import (
    "bufio"
    "io"
    "log"
    "net/http"
    "os"
    "sync/atomic"

    "github.com/Tomilla/imagespider/config"
    "github.com/Tomilla/imagespider/net"
    "github.com/Tomilla/imagespider/util"
)

type work struct {
    url      string
    fileName string
}

func newWork(url string, fileName string) work {
    return work{url: url, fileName: fileName}
}

type worker struct {
    s *scheduler
}

func NewWorkers(s *scheduler) *worker {
    return &worker{s: s}
}

func (w *worker) Start() {

    for i := 0; i < w.s.workerCount; i++ {
        go w.work()
    }

}
func (w *worker) work() {
    workChan := make(chan work)
    w.s.Ready(workChan)

    for {
        task, ok := <-workChan
        if !ok {
            return // channel 关闭，退出
        }

        w.Download(task)
        util.SleepRandomDuration(config.C.GetSleepRange())
        w.s.Ready(workChan)

    }

}

func (w *worker) Download(task work) {

    err := w.downloadWithPath(task.url, task.fileName)
    if err != nil {
        log.Println("####### Error download ", err, task.url)
        util.WarnErr(os.Remove(task.fileName)) // 下载失败 删除文件
        return
    }

    atomic.AddInt32(&count, 1)

    // log.Printf("#%d downloaded %s", count, task.fileName)

}

func (w *worker) downloadWithPath(link, fileName string) error {

    if util.CheckPathExists(fileName) {
        return nil
    }
    // resp, err := http.Get(link)
    // @@@@@@@@@@@@@@@@@

    client := net.NewClient(false)
    req, err := http.NewRequest("GET", link, nil)
    if err != nil {
        return nil
    }
    // req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36 t66y_com")
    req.Header.Set("User-Agent", net.GetRandomUserAgent())
    resp, err := client.Do(req)

    if err != nil {
        return err
    }
    defer func() {
        util.WarnErr(resp.Body.Close())
    }()
    buf := bufio.NewReader(resp.Body)
    newFile, err := os.Create(fileName)
    if err != nil {
        return err
    }

    _, err = io.Copy(newFile, buf)
    util.WarnErr(err)
    defer func() {
        util.WarnErr(newFile.Close())
    }()
    return nil
}

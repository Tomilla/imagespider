package image

import (
    "bufio"
    "io"
    "log"
    "net/http"
    "os"
    "path"
    "sync/atomic"

    "github.com/Tomilla/imagespider/config"
    "github.com/Tomilla/imagespider/glog"
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

    if glog.CheckPathExists(fileName) {
        return nil
    }
    client := net.NewClient(false)
    req, err := http.NewRequest("GET", link, nil)
    if err != nil {
        return nil
    }
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
        if config.C.GetShowDownloadProgress() {
            config.L.Infof("Image Downloaded: %v\n%v\n", path.Base(fileName), link)
        }
    }()
    return nil
}

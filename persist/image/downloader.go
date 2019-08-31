package image

import (
    "fmt"
    netUrl "net/url"
    "os"
    "path"
    "strconv"
    "strings"
    "unicode/utf8"

    "github.com/google/uuid"

    "github.com/Tomilla/imagespider/collections/set"
    "github.com/Tomilla/imagespider/common"
    "github.com/Tomilla/imagespider/util"
)

var (
    ValidImageExtension = set.New(".jpg", ".jpeg", ".gif", ".png")
)

const (
    DefaultImageExtension = ".jpg"
)
type downloader struct {
    common.ImageConfig
    common.Limit
    workChan chan work
}

func NewDownloader(imageConfig *common.ImageConfig, limit *common.Limit) *downloader {
    return &downloader{ImageConfig: *imageConfig, Limit: *limit}
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
            fileName := d.GetFileName(baseFolder, topic.Name, url, i)
            w := newWork(topic.Key, url, fileName)
            workChan <- w
        }

    }
}

func (d *downloader) CreateWorker(s *scheduler) {
    ws := NewWorkers(s)
    ws.Start()
}

func (d *downloader) GetFileName(baseFolder, postName string, postUrl string, imgIndex int) string {
    limit := d.ImagePathLenLimit
    extension := strings.ToLower(path.Ext(postUrl))
    if len(extension) == 0 || !ValidImageExtension.Has(extension) {
        extension = DefaultImageExtension
    }
    if d.UniqFolder { // only provide filename is fine for unique folder
        return baseFolder +
            strings.ToLower("_"+util.LeftPad2Len(strconv.Itoa(imgIndex), "0", 3)+extension)
    }
    // otherwise, we need join the post postName with filename
    u, err := netUrl.Parse(postUrl)
    if err != nil {
        postName = uuid.New().String()
    } else {
        postName = strings.ToLower(u.Path)
        postName = path.Base(postName)
        lenCharacters := utf8.RuneCountInString(postName)
        if lenCharacters >= limit {
            r := []rune(postName)
            r = r[lenCharacters-limit:]
            postName = string(r)
        }
    }
    return path.Join(baseFolder, fmt.Sprintf("%03v_%v", imgIndex, postName))
}

package parser

import (
    "bytes"
    "fmt"
    "io/ioutil"
    netUrl "net/url"
    "os"
    "path"
    "strings"

    "github.com/PuerkitoBio/goquery"

    "github.com/Tomilla/imagespider/common"
    "github.com/Tomilla/imagespider/common/model"
    "github.com/Tomilla/imagespider/engine"
    "github.com/Tomilla/imagespider/glog"
    "github.com/Tomilla/imagespider/util"
)

type PostRequest struct {
    URL   string
    Agent string
    Post  *engine.Post
}

func (p PostRequest) GetURL() string {
    return p.URL
}

func (p *PostRequest) SetURL(new string) bool {
    p.URL = new
    return true
}

func (p PostRequest) GetAgent() string {
    return p.Agent
}

func (p PostRequest) GetPost() *engine.Post {
    return p.Post
}

func (p PostRequest) Parser(contents []byte, url string) *engine.ParseResult {
    if !p.Archiver(contents, url) {
        common.L.Infof("Cannot archive content of :", url)
    }
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
    if err != nil {
        common.Redis.HSet(SimplifyPostUrl(url), common.TopicEnum.Status, common.PostContentFailParsed)
        return nil
    }
    node := doc.Find(".tpc_content.do_not_catch").ParentsUntil("table")
    if len(node.Nodes) == 0 {
        node = doc.Find("h4").ParentsUntil("table")
        if len(node.Nodes) == 0 {
            return nil
        }
    }
    // html, err := node.Html()
    // if err != nil {
    //     return nil
    // }
    // println(html)
    matches := make(chan string)
    go ExtractImageUrls(node, matches)

    // imageMatches := imageRe.FindAllSubmatch(contents, -1)
    // if imageMatches == nil {
    //     log.Println("nil images")
    //     return nil
    // }

    t := model.Topic{}
    key := SimplifyPostUrl(url)
    t.Url = url
    t.Name = p.Post.Title
    t.CountImage = p.Post.CountImage
    t.CountReply = p.Post.CountReply
    t.Key = key
    common.Redis.HSet(key, common.TopicEnum.CountImage, t.CountImage)
    common.Redis.HSet(key, common.TopicEnum.CountReply, t.CountReply)
    common.Redis.HSet(key, common.TopicEnum.Url, t.Url)
    common.Redis.HSet(key, common.TopicEnum.Name, t.Name)

    _name := p.Post.Title
    t.Name = fmt.Sprintf(FILE_NAME_FORMAT, t.CountReply, t.CountImage, NormalizeName(_name))
    t.Url = url

    for m := range matches {
        if !isInvalid(m) {
            t.Images = append(t.Images, m)
        }
    }
    // the post is still alive
    if len(t.Images) >= t.CountImage {
        common.Redis.HSet(key, common.TopicEnum.Status, common.PostImgAllParsed)
    } else if len(t.Images) > 0 {
        common.Redis.HSet(key, common.TopicEnum.Status, common.PostImgPartParsed)
    } else {
        common.Redis.HSet(key, common.TopicEnum.Status, common.PostImgFailParsed)
    }
    return &engine.ParseResult{Items: []interface{}{t}}
}

func (p PostRequest) Archiver(contents []byte, url string) bool {

    logPath := common.C.GetLogPath()
    if !glog.CheckPathExists(logPath) {
        err := os.MkdirAll(logPath, util.DefaultFilePerm)
        if err != nil {
            common.L.Debug("Cannot create logPath")
            return false
        }
    }
    u, err := netUrl.Parse(url)
    if err != nil {
        common.L.Debug("Cannot parse url")
        return false
    }
    err = ioutil.WriteFile(
        path.Join(logPath, NormalizePostUrl(u.Path, true)),
        contents,
        0664)
    if err != nil {
        common.L.Error("Cannot write logPath")
        return false
    }
    return true
}

// delete duplicates
func isInvalid(s string) bool {
    result := false
    switch {
    // exclude non-image
    // such as: https://www.yuoimg.com/u/20190218/12543160.jpg
    case strings.Contains(s, `/i/?i=u`):
        result = true
        break
    case strings.Contains(s, `imagexport`):
        result = true
        break
    default:
        break
    }
    return result
}

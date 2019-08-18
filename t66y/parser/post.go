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

    "github.com/Tomilla/imagespider/config"
    "github.com/Tomilla/imagespider/engine"
    "github.com/Tomilla/imagespider/glog"
    "github.com/Tomilla/imagespider/model"
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
        config.L.Infof("Cannot archive content of :", url)
    }
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
    if err != nil {
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
    t.CountImage = p.Post.CountImage
    t.CountReply = p.Post.CountReply

    _name := p.Post.Title
    t.Name = fmt.Sprintf("[%03v][%03v]%v", t.CountReply, t.CountImage, NormalizeName(_name))
    t.Url = url

    for m := range matches {
        if !isDup(m) {
            t.Images = append(t.Images, m)
        }
    }
    return &engine.ParseResult{Items: []interface{}{t}}
}

func (p PostRequest) Archiver(contents []byte, url string) bool {

    logPath := config.C.GetLogPath()
    if !glog.CheckPathExists(logPath) {
        err := os.MkdirAll(logPath, util.DefaultFilePerm)
        if err != nil {
            config.L.Debug("Cannot create logPath")
            return false
        }
    }
    u, err := netUrl.Parse(url)
    if err != nil {
        config.L.Debug("Cannot parse url")
        return false
    }
    err = ioutil.WriteFile(
        path.Join(logPath, strings.Trim(postPathRe.ReplaceAllString(u.Path, "_"), "_")),
        contents,
        0664)
    if err != nil {
        config.L.Error("Cannot write logPath")
        return false
    }
    return true
}

// delete duplicates
func isDup(s string) bool {
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

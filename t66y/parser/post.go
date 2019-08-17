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

func (p PostRequest) SetURL(new string) {
    p.URL = new
}

func (p PostRequest) GetAgent() string {
    return p.Agent
}

func (p PostRequest) GetPost() *engine.Post {
    return p.Post
}

func (p PostRequest) Parser(contents []byte, url string) *engine.ParseResult {
    logPath := config.C.GetLogPath()
    if !glog.CheckPathExists(logPath) {
        err := os.MkdirAll(logPath, util.DefaultFilePerm)
        if err != nil {
            config.L.Debug("Cannot create logPath")
        }
    }
    u, err := netUrl.Parse(url)
    if err != nil {
        config.L.Debug("Cannot parse url")
        return nil
    }
    err = ioutil.WriteFile(
        path.Join(logPath, PathSeparatorRe.ReplaceAllString(strings.Trim(u.Path, `\/`), "_")),
        contents,
        0664)
    if err != nil {
        config.L.Error("Cannot write logPath")
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

func NormalizeName(s string) string {
    // s = strings.Trim(s, "[]")
    // fmt.Println("before -- > ", s)
    limit := config.C.GetPostNameLenLimit()
    // remove quote
    s = quoteRe.ReplaceAllString(s, `$1`)
    // remove punctuation
    s = punctuationRe.ReplaceAllString(s, `$1`)
    // replace full characters to half characters
    s = halfToFullRe.ReplaceAllString(s, `[$1]`)
    if strings.Contains(s, `/`) { // 去除名字中的反斜杠
        s = strings.Replace(s, `/`, ``, -1)
    }
    // tags := tagBracketRe.FindAllString(s, -1)
    // fmt.Printf("tags: %v\n", tags)
    s = tagBracketRe.ReplaceAllString(s, "")
    s = nonTagBracketRe.ReplaceAllString(s, "$1")
    // replace non Chinese or non English world
    s = nonChineseEnglishRe.ReplaceAllString(s, " ")
    s = whiteSpaceInsideRe.ReplaceAllString(s, "$1 _ $2")
    s = whiteSpaceRe.ReplaceAllString(s, "")

    characters := []rune(s)

    if len(characters) > limit {
        characters = characters[:limit]
    }
    return string(characters)
}

func ExtractImageUrls(content interface{}, matches chan string) {
    switch content.(type) {
    case *goquery.Selection:
        content := content.(*goquery.Selection)
        found := content.Find("input[type=image]")
        for _, node := range found.Nodes {
            inputTag := goquery.NewDocumentFromNode(node)
            val, exist := inputTag.Attr("data-src")
            if !exist {
                val, exist = inputTag.Attr("src")
                if !exist {
                    html, err := inputTag.Html()
                    if err != nil {
                        config.L.Debug("It's not a valid image tag.")
                    } else {
                        config.L.Debugf("It's not a valid image tag.\n", html)
                    }
                }
            }
            if viidiiImageRe.MatchString(val) {
                continue
            }
            if postImageRe.MatchString(val) {
                matches <- val
            }
        }
        break
    case string:
        content := content.(string)
        doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
        if err != nil {
            config.L.Debugf("Cannot create document from content: %v\n", content)
        }
        ExtractImageUrls(doc, matches)
        break
    case []byte:
        content := content.([]byte)
        doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
        if err != nil {
            config.L.Debugf("Cannot create document from content: %v\n", content)
        }
        ExtractImageUrls(doc, matches)
    default:
        break
    }
    close(matches)
}

package parser

import (
    "bytes"
    "fmt"
    "io/ioutil"
    netUrl "net/url"
    "os"
    "path"
    "strconv"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/corpix/uarand"

    "github.com/Tomilla/imagespider/common"
    "github.com/Tomilla/imagespider/engine"
    "github.com/Tomilla/imagespider/glog"
    "github.com/Tomilla/imagespider/util"
)

type PostListRequest struct {
    URL   string
    Agent string
}

func (p PostListRequest) GetURL() string {
    return p.URL
}

func (p *PostListRequest) SetURL(new string) bool {
    p.URL = new
    return true
}

func (p PostListRequest) GetAgent() string {
    return p.Agent
}

func (p PostListRequest) GetPost() *engine.Post {
    return nil
}

func (p PostListRequest) Parser(contents []byte, url string) *engine.ParseResult {
    if !p.Archiver(contents, url) {
        common.L.Infof("Cannot archive content of %v", url)
    }
    fmt.Println(url)
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
    util.CheckErr(err)
    allTableRow := doc.Find("tr.tr3.t_one.tac")

    replyLow, replyHigh := common.C.GetReplyRange() // limit the topic count
    result := engine.ParseResult{}

    var postColor string
    allTableRow.Each(func(n int, sel *goquery.Selection) {
        // baseHtml, err := sel.Html()

        allTd := sel.Find("td")
        /*
         * 0: post_href(duplicated, see below)
         * 1: post_title, post_href
         * 2: post_author, post_time
         * 3: reply_count
         * 4: reply_last_time, reply_last_author
         */
        if len(allTd.Nodes) < 5 {
            return
        }

        // fmt.Println(util.LeftPad2Len("", "-", 80))
        // fmt.Printf("%05v %v\n", n, baseHtml)
        // fmt.Println(util.RightPad2Len("", "-", 80))
        db := common.DB
        if db != nil {
            // do something
        }
        var post = engine.Post{}

        for i, node := range allTd.Nodes {
            doc := goquery.NewDocumentFromNode(node)
            switch i {
            case 0:
                // fmt.Println(util.LeftPad2Len("", "*", 80))
                break
            case 1:
                var exist bool
                aTag := doc.Find("h3>a")
                _ref, exist := aTag.Attr("href")
                if exist {
                    if postUrlRe.MatchString(_ref) {
                        post.Path = _ref
                    } else {
                        return // ignore invalid url
                    }
                }

                _title := aTag.Text()
                post.Title = _title
                postColor, exist = aTag.Find("font").Attr("color")
                if exist && ignoredPostColor.Has(postColor) {
                    fmt.Printf("Ignore Admin Post: %v %v\n", post.Title, postColor)
                    return
                }

                _imageCount, err := strconv.ParseInt(
                    util.GetLastItem(countImageRe.FindAllStringSubmatch(_title, -1)),
                    10, 64)
                if err != nil {
                    fmt.Println("cannot parse image count")
                    post.CountImage = 0
                } else {
                    post.CountImage = int(_imageCount)
                }
                break
            case 2:
                post.Author = doc.Find("a").Text()
                _createdAt := strings.Trim(doc.Find("div").Text(), util.WhiteSpace)
                post.CreatedAt, err = time.Parse(util.DateLayout, _createdAt)
                if err != nil {
                    post.CreatedAt, _ = time.Parse(util.DateLayout, util.DateDefault)
                }
                post.CreatedAt = post.CreatedAt.In(util.TZ)
                break
            case 3:
                _replyCount, err := strconv.ParseInt(strings.Trim(doc.Text(), util.WhiteSpace), 10, 64)
                if err != nil {
                    continue
                }
                _replyCountInt := int(_replyCount)

                if !includePostColor.Has(postColor) && !(_replyCountInt >= replyLow && _replyCountInt < replyHigh) {
                    fmt.Printf("Ignore Post: %v %v\n", post.Path, _replyCount)
                    return
                }
                post.CountReply = _replyCountInt
            case 4:
                _lastReplyAt := strings.Trim(doc.Find("a").Text(), util.WhiteSpace)
                post.LastReplyAt, err = time.Parse(util.DateTimeLayout, _lastReplyAt)
                if err != nil {
                    post.LastReplyAt, _ = time.Parse(util.DateTimeLayout, util.DateTimeDefault)
                }
                post.LastReplyAt = post.LastReplyAt.In(util.TZ)
            default:
                s, err := doc.Html()
                if err != nil {
                    continue
                }
                fmt.Println(i, s)
            }
        }

        // fmt.Println(pTitle)
        fmt.Println(post)
        result.Items = append(result.Items, "topic: "+post.Title)
        result.Requests = append(result.Requests, &PostRequest{
            URL:   HOSTNAME + post.Path,
            Agent: uarand.GetRandom(),
            Post:  &post,
        })
    })
    return &result
}

func (p PostListRequest) Archiver(contents []byte, url string) bool {
    u, err := netUrl.Parse(url)
    if err != nil {
        common.L.Infof("Cannot parse url %v: %v", url, err)
        return false
    }

    logPath := common.C.GetLogPath()
    if !glog.CheckPathExists(logPath) {
        err := os.MkdirAll(logPath, util.DefaultFilePerm)
        if err != nil {
            common.L.Infof("Cannot MkdirAll for %v: %v", logPath, err)
            return false
        }
    }

    uBase := path.Base(u.Path)
    uExt := path.Ext(u.Path)

    if len(uExt) > 0 {
        uBase = strings.ReplaceAll(uBase, uExt, "")
    }
    if !ValidWebPageExt.Has(uExt) {
        uExt = DefaultWebPageExt
    }

    finalPath := strings.ToLower(uBase + "_" + u.RawQuery + uExt)
    err = ioutil.WriteFile(
        path.Join(logPath, strings.Trim(postPathRe.ReplaceAllString(finalPath, "_"), "_")),
        contents,
        0664)
    if err != nil {
        common.L.Infof("Cannot WriteFile of %v: %v", finalPath, err)
        return false
    }
    return true
}

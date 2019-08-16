package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	netUrl "net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"

	"github.com/Tomilla/imagespider/collections/set"
	"github.com/Tomilla/imagespider/config"
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

func (p PostListRequest) SetURL(new string) {
	p.URL = new
}

func (p PostListRequest) GetAgent() string {
	return p.Agent
}

func (p PostListRequest) GetPost() *engine.Post {
	return nil
}

var (
	postUrlRe        = regexp.MustCompile(`htm_data/\d+/\d+/\d+\.html`)
	countImageRe     = regexp.MustCompile(`[\[［](\d+).?[］\]]`)
	ignoredPostColor = set.New("red", "blue", "orange")
	includePostColor = set.New("green")
)

const (
	HOSTNAME = "http://t66y.com/"
)

func (p PostListRequest) Parser(contents []byte, url string) engine.ParseResult {
	fmt.Println(url)
	fmt.Println(util.LeftPad2Len("", "*", 80))
	// fmt.Println(string(contents))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	util.CheckErr(err)
	allTableRow := doc.Find("tr.tr3.t_one.tac")

	replyLow, replyHigh := config.C.GetReplyRange() // limit the topic count
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
		db := config.DB
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
		result.Items = append(result.Items, "topic: "+string(post.Title))
		result.Requests = append(result.Requests, PostRequest{
			URL:   HOSTNAME + string(post.Path),
			Agent: uarand.GetRandom(),
			Post:  &post,
		})
	})
	u, err := netUrl.Parse(url)
	util.CheckErr(err)
	logPath := config.C.GetLogPath()
	if !glog.CheckPathExists(logPath) {
		err := os.MkdirAll(logPath, util.DefaultFilePerm)
		util.CheckErr(err)
	}

	err = ioutil.WriteFile(
		path.Join(logPath, strings.Replace(u.Path, string(os.PathSeparator), "_", -1)),
		contents,
		0664)
	util.CheckErr(err)

	fmt.Println(util.RightPad2Len("", "*", 80))
	return result
}
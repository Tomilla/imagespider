package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"

	"github.com/Tomilla/imagespider/config"
	"github.com/Tomilla/imagespider/engine"
	"github.com/Tomilla/imagespider/util"
)

var topicListRe = regexp.MustCompile(`<h3><a href="(htm_data/[0-9]*/[0-9]*/[0-9]*\.html)"[^>]*>([^<]+)</a>`)

func ParseTopicList(contents []byte, url string) engine.ParseResult {
	fmt.Println(url)
	fmt.Println(util.LeftPad2Len("", "*", 80))
	// fmt.Println(string(contents))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	util.CheckErr(err)
	allTableRow := doc.Find("tr.tr3.t_one.tac")

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
		}
		var (
			ref         string
			author      string
			createdAt   time.Time
			lastReplyAt time.Time
			replyCount  int64
		)

		for i, node := range allTd.Nodes {
			doc := goquery.NewDocumentFromNode(node)
			switch i {
			case 0:
				break
			case 1:
				_ref, exist := doc.Find("h3>a").Attr("href")
				if exist {
					fmt.Println(ref)
					ref = _ref
				}
				break
			case 2:
				author = doc.Find("a").Text()
				_createdAt := strings.Trim(doc.Find("div").Text(), util.WhiteSpace)
				createdAt, err = time.Parse(util.DateLayout, _createdAt)
				if err != nil {
					createdAt, _ = time.Parse(util.DateLayout, util.DateDefault)
				}
				fmt.Println(author, createdAt.In(util.TZ))
				break
			case 3:
				replyCount, err = strconv.ParseInt(strings.Trim(doc.Text(), util.WhiteSpace), 10, 64)
				fmt.Println(replyCount)
			case 4:
				_lastReplyAt := strings.Trim(doc.Find("a").Text(), util.WhiteSpace)
				lastReplyAt, err = time.Parse(util.DateTimeLayout, _lastReplyAt)
				if err != nil {
					lastReplyAt, _ = time.Parse(util.DateTimeLayout, util.DateTimeDefault)
				}
				fmt.Println(lastReplyAt.In(util.TZ))
			default:
				s, err := doc.Html()
				if err != nil {
					continue
				}
				fmt.Println(i, s)
			}
		}

		pTitle := sel.Find("h3>a").Text()

		fmt.Println(pTitle)
	})

	err = ioutil.WriteFile(config.C.GetLogPath(), contents, 0664)
	util.CheckErr(err)

	fmt.Println(util.RightPad2Len("", "*", 80))

	matches := topicListRe.FindAllSubmatch(contents, -1)
	limit := config.C.GetPageLimit() // limit the topic count
	result := engine.ParseResult{}
	for _, m := range matches {
		result.Items = append(result.Items, "topic: "+string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:        "http://t66y.com/" + string(m[1]),
			Agent:      uarand.GetRandom(),
			ParserFunc: ParseTopic,
			Name:       string(m[2]),
		})
		limit--
		if limit < 0 {
			return result
		}

	}
	return result

}

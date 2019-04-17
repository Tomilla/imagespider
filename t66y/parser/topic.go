package parser

import (
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/engine"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"
	"log"
	"regexp"
)

//var imageRe = regexp.MustCompile(`src=(http(s)?://[^/]*/u/[0-9]+/[0-9]+\.jpg)'`)
//var imageRe = regexp.MustCompile(`['"](http(s)?://[^'"]+jpg)["']`)
var imageRe = regexp.MustCompile(`data-(src|link)=['"](http[s]?://[^'"]+)['"]`)
var titleRe = regexp.MustCompile(`<title>([^<]+)</title>`)
var ImageCh = make(chan []*model.Topic, 20)

func ParseTopic(contents []byte) engine.ParseResult {

	log.Printf("ParseTopic ")

	imageMatches := imageRe.FindAllSubmatch(contents, -1)
	if imageMatches == nil {
		panic("nil images")
	}

	titleMatch := titleRe.FindSubmatch(contents)

	t := &model.Topic{}
	name := string(titleMatch[1])

	t.Name = normalizeName(name)
	log.Println("["+t.Name+"]"+" matches -->  ", len(imageMatches), "images")

	for _, m := range imageMatches {
		t.Images = append(t.Images, string(m[2]))
	}
	go func(topic *model.Topic) { // send to download channel
		ImageCh <- []*model.Topic{t}
	}(t)
	return engine.ParseResult{}
}

func normalizeName(s string) string {
	r := []rune(s)
	for i, v := range r {
		if v == ' ' {
			return string(r[:i])
		}

	}
	return ""
}

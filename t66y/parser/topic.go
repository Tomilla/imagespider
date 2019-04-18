package parser

import (
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/engine"
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"
	"regexp"
)

var imageRe = regexp.MustCompile(`data-(src|link)=['"](http[s]?://[^'"]+)['"]`)
var titleRe = regexp.MustCompile(`<title>([^<]+)</title>`)
var ImageCh = make(chan []*model.Topic, 20)

func ParseTopic(contents []byte) engine.ParseResult {

	imageMatches := imageRe.FindAllSubmatch(contents, -1)
	if imageMatches == nil {
		panic("nil images")
	}

	titleMatch := titleRe.FindSubmatch(contents)

	t := model.Topic{}
	name := string(titleMatch[1])

	t.Name = normalizeName(name)
	//log.Println("ParseTopic ["+t.Name+"]"+" matches -->  ", len(imageMatches), "images")

	for _, m := range imageMatches {
		t.Images = append(t.Images, string(m[2]))
	}

	return engine.ParseResult{Items: []interface{}{t}}
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

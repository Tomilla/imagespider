package parser

import (
	"regexp"
	"strings"

	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/engine"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"
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

	for _, m := range imageMatches {
		//filter(m[2])
		if isDup(m[2]) {
			continue
		}
		t.Images = append(t.Images, string(m[2]))
	}

	return engine.ParseResult{Items: []interface{}{t}}
}

func normalizeName(s string) string {
	//s = strings.Trim(s, "[]")
	//fmt.Println("before -- > ", s)
	r := []rune(s)
	for i, v := range r {
		if v == '-' {
			result := string(r[:i-1])
			//fmt.Println("after --> ", result)

			return result
		}
	}
	return ""
}

// delete duplicates
func isDup(b []byte) bool {
	s := string(b)
	switch {
	case strings.Contains(s, `/i/?i=u`): //并不是图片文件
		return true

	default:
		return false

	}
}

func filter(b []byte) []byte {

	if !isDup(b) {
		return b
	}

	s := string(b)
	switch {
	case strings.Contains(s, `/i/?i=u`):
		//fmt.Println("before  Replaced ", s)
		s := strings.Replace(s, `i/?i=u`, `u`, -1)

		//fmt.Println("after Replaced ", s)
		return []byte(s)
	default:
		return b
	}

}

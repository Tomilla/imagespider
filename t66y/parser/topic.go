package parser

import (
	"log"
	"regexp"
	"strings"

	"github.com/Tomilla/imagespider/config"

	"github.com/Tomilla/imagespider/engine"
	"github.com/Tomilla/imagespider/model"
)

// var imageRe = regexp.MustCompile(`(data-src|data-link|src)=['"](http[s]?://[^'"]+[^s])['"]`)

var imageRe = regexp.MustCompile(`(?i)(data-src|data-link|src)=['"](http[s]?://[^'"]+(jpg|png|jpeg|gif))['"]`)
var titleRe = regexp.MustCompile(`<title>([^>]+)(\s+-\s*\S+\s*\|\s*\S+\s*-\s*\S+\s*)</title>`)

// var ImageCh = make(chan []*model.Topic, 20)

func ParseTopic(contents []byte, url string) engine.ParseResult {

	imageMatches := imageRe.FindAllSubmatch(contents, -1)
	if imageMatches == nil {
		log.Println("nil images")
		return engine.ParseResult{}
	}

	titleMatch := titleRe.FindSubmatch(contents)

	t := model.Topic{}
	name := string(titleMatch[1])

	println("-->", name)
	t.Name = normalizeName(name)
	println("==>", t.Name)
	t.Url = url

	for _, m := range imageMatches {
		url := string(m[2])
		if isDup(url) {
			continue
		}
		t.Images = append(t.Images, string(m[2]))
		// fmt.Println("added", string(m[2]), t.Name)
	}

	return engine.ParseResult{Items: []interface{}{t}}
}

func normalizeName(s string) string {
	// s = strings.Trim(s, "[]")
	// fmt.Println("before -- > ", s)
	limit := config.C.GetNameLenLimit()
	result := s

	if strings.Contains(result, `/`) { // 去除名字中的反斜杠
		result = strings.Replace(result, `/`, ``, -1)
	}
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

// delete duplicates
func isDup(s string) bool {
	result := false
	switch {
	case strings.Contains(s, `/i/?i=u`): // 并不是图片文件
		result = true // 如 https://www.yuoimg.com/u/20190218/12543160.jpg
	// case strings.Contains(s, `www.kanjiantu.com/image/`):
	// 	result = true
	// case strings.Contains(s, `sb88y.net`):
	// 	result = true
	// case strings.Contains(s, `htm`):
	// 	result = true
	// case strings.Contains(s, `h34229`):
	// 	result = true
	// case strings.Contains(s, `img599`):
	// 	result = true
	// case strings.Contains(s, `667um`):
	// 	result = true
	// case strings.Contains(s, `51668`):
	// 	result = true
	// case strings.Contains(s, `x6img`):
	// 	result = true
	// case strings.Contains(s, `dioimg`):
	// 	result = true
	// case strings.Contains(s, `sinaimg`):
	//  result = true
	case strings.Contains(s, `imagexport`): // 这个网址不直接提供图片文件
		result = true
	// case strings.Contains(s, `?`):
	// 	return true
	default:

	}
	return result
}

// func filter(b []byte) []byte {
//
// 	if !isDup(b) {
// 		return b
// 	}
//
// 	s := string(b)
// 	switch {
// 	case strings.Contains(s, `/i/?i=u`):
// 		//fmt.Println("before  Replaced ", s)
// 		s := strings.Replace(s, `i/?i=u`, `u`, -1)
//
// 		//fmt.Println("after Replaced ", s)
// 		return []byte(s)
// 	default:
// 		return b
// 	}
//
// }

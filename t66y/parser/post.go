package parser

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Tomilla/imagespider/config"

	"github.com/Tomilla/imagespider/engine"
	"github.com/Tomilla/imagespider/model"
)

var (
	imageRe = regexp.MustCompile(`(?i)(data-src|data-link|src)=['"](http[s]?://[^'"]+(jpg|png|jpeg|gif))['"]`)
	titleRe = regexp.MustCompile(`<title>([^>]+)(\s+-\s*\S+\s*\|\s*\S+\s*-\s*\S+\s*)</title>`)
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

func (p PostRequest) Parser(contents []byte, url string) engine.ParseResult {
	imageMatches := imageRe.FindAllSubmatch(contents, -1)
	if imageMatches == nil {
		log.Println("nil images")
		return engine.ParseResult{}
	}

	titleMatch := titleRe.FindSubmatch(contents)

	t := model.Topic{}
	name := string(titleMatch[1])
	t.CountImage = p.Post.CountImage
	t.CountReply = p.Post.CountReply

	println("-->", name)
	t.Name = fmt.Sprintf("[%v][%v]%v", t.CountReply, t.CountImage, normalizeName(name))
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
	if strings.Contains(s, `/`) { // 去除名字中的反斜杠
		s = strings.Replace(s, `/`, ``, -1)
	}

	characters := []rune(s)

	if len(characters) > limit {
		characters = characters[:limit]
	}
	return string(characters)
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

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
	// template
	tmpChineseEnglish   = `[%v\p{Han}\p{Latin}0-9_-]`
	chineseEnglish      = fmt.Sprintf(tmpChineseEnglish, ``)
	nonChineseEnglish   = fmt.Sprintf(tmpChineseEnglish, `^`)
	imageRe             = regexp.MustCompile(`(?i)(data-src|data-link|src)=['"](http[s]?://[^'"]+(jpg|png|jpeg|gif))['"]`)
	titleRe             = regexp.MustCompile(`<title>([^>]+)(\s+-\s*\S+\s*\|\s*\S+\s*-\s*\S+\s*)</title>`)
	quoteRe             = regexp.MustCompile(`(?:['"‘“])(.*?)(?:['"’”])`)
	punctuationRe       = regexp.MustCompile(`(?:[(（{｛])(.*?)(?:\[\)）}｝])`)
	halfToFullRe        = regexp.MustCompile(`(?:[［「【『〖])(.*?)(?:[］」】』〗])`)
	tagBracketRe        = regexp.MustCompile(`\[.{0,6}]`)
	nonTagBracketRe     = regexp.MustCompile(`\[(.*?)]`)
	whiteSpaceRe        = regexp.MustCompile(`\s+`)
	whiteSpaceInsideRe  = regexp.MustCompile(fmt.Sprintf(`(%v)\s+(%v)`, chineseEnglish, chineseEnglish))
	nonChineseEnglishRe = regexp.MustCompile(fmt.Sprintf(`%v+`, nonChineseEnglish))
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
	t.Name = fmt.Sprintf("[%03v][%03v]%v", t.CountReply, t.CountImage, NormalizeName(name))
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
	limit := config.C.GetNameLenLimit()
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

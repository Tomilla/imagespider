package parser

import (
	"fmt"

	"github.com/corpix/uarand"

	"github.com/wuxiangzhou2010/imagespider/config"
	"github.com/wuxiangzhou2010/imagespider/engine"

	"regexp"
)

var topicListRe = regexp.MustCompile(`<h3><a href="(htm_data/[0-9]*/[0-9]*/[0-9]*\.html)"[^>]*>([^<]+)</a>`)

func ParseTopicList(contents []byte, url string) engine.ParseResult {
	fmt.Println(url)
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

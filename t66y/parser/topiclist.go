package parser

import (
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler/util/agent/my"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/config"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/engine"

	"regexp"
)

var topicListRe = regexp.MustCompile(`<h3><a href="(htm_data/[0-9]*/[0-9]*/[0-9]*\.html)"[^>]*>([^<]+)</a>`)

func ParseTopicList(contents []byte) engine.ParseResult {

	matches := topicListRe.FindAllSubmatch(contents, -1)
	limit := config.C.GetPageLimit() // limit the topic count
	result := engine.ParseResult{}
	for _, m := range matches {
		result.Items = append(result.Items, "topic: "+string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:        "http://t66y.com/" + string(m[1]),
			Agent:      my.NewAgent(),
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

package parser

import (
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/engine"

	"regexp"
)

//var cityListRe = regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
//var cityListRe = regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
var topicListRe = regexp.MustCompile(`<h3><a href="(htm_data/[0-9]*/[0-9]*/[0-9]*\.html)"[^>]*>([^<]+)</a>`)

func ParseTopicList(contents []byte) engine.ParseResult {

	//re := regexp.MustCompile(cityListRe)
	matches := topicListRe.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, m := range matches {
		result.Items = append(result.Items, "City: "+string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url: "http://t66y.com/" + string(m[1]),
			//ParserFunc: ParseTopic,
			ParserFunc: ParseTopic,
		})

	}
	return result

}

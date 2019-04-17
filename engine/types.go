package engine

import "github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler/util/agent"

type Request struct {
	Url        string
	ParserFunc func([]byte) ParseResult
	Agent      agent.Agent
	Name       string
}
type ParseResult struct {
	Requests []Request
	Items    []interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}

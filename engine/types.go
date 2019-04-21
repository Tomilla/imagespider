package engine

import "github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler/util/agent"

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

func NewParseResult(items []interface{}) *ParseResult {
	return &ParseResult{Items: items}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}

type Engine interface {
	Run(s Scheduler, request []Request)
	Shutdown()
}

type Scheduler interface {
	Schedule(chan bool)
	SubmitRequest(Request)
	SubmitWorker(chan Request)
	GetWorkCount() int
	Shutdown()
}

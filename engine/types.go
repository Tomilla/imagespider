package engine

import (
	"time"
)

type Request struct {
	Url        string
	ParserFunc func([]byte, string) ParseResult
	Agent      string
	Name       string
}
type ParseResult struct {
	Requests []Request
	Items    []interface{}
}

type Post struct {
	CreatedAt   time.Time
	LastReplyAt time.Time
	Author      string
	Path        string
	Title       string
	ReplyCount  int64
}

// func NewParseResult(items []interface{}) *ParseResult {
//     return &ParseResult{Items: items}
// }

// func NilParser([]byte) ParseResult {
//     return ParseResult{}
// }

type Engine interface {
	Run(s Scheduler, request []Request)
	Shutdown()
}

type Scheduler interface {
	Schedule(chan bool)
	SubmitRequest(Request)
	SubmitWorker(chan Request)

	Shutdown()
}

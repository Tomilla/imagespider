package engine

import (
    "time"
)

type BaseParser interface {
    Parser([]byte, string) *ParseResult
    GetURL() string
    SetURL(string)
    GetAgent() string
    GetPost() *Post
}

type Request struct {
    URL        string
    Agent      string
    Post       *Post
    ParserFunc func([]byte, string) ParseResult
}

type ParseResult struct {
    Requests []BaseParser
    Items    []interface{}
}

type Post struct {
    CreatedAt   time.Time
    LastReplyAt time.Time
    Author      string
    Path        string
    Title       string
    CountReply  int
    CountImage  int
}

// func NewParseResult(items []interface{}) *ParseResult {
//     return &ParseResult{Items: items}
// }

// func NilParser([]byte) ParseResult {
//     return ParseResult{}
// }

type Engine interface {
    Run(s Scheduler, request []BaseParser)
    Shutdown()
}

type Scheduler interface {
    Schedule(chan bool)
    SubmitRequest(parser BaseParser)
    SubmitWorker(chan BaseParser)

    Shutdown()
}

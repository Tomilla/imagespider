package generator

import (
	"strconv"

	"github.com/corpix/uarand"

	"github.com/Tomilla/imagespider/config"
	"github.com/Tomilla/imagespider/engine"
	"github.com/Tomilla/imagespider/t66y/parser"
)

const nextString = "&search=&page="

type Generator struct {
	seeds         []string
	startRequests []engine.Request
	count         int
	endPageNum    int
	requestChan   chan engine.Request
}

func NewGenerator(seeds []string) chan engine.Request {
	g := &Generator{
		seeds:      seeds,
		count:      config.C.GetStartPageNum(),
		endPageNum: config.C.GetEndPageNum(),
	}
	g.requestChan = make(chan engine.Request)
	go g.Generate()
	return g.requestChan
}

func (g *Generator) SetRequestChan(requestChan chan engine.Request) {
	g.requestChan = requestChan

}
func (g *Generator) Generate() {
	g.GenerateStartRequest(g.seeds)
	for {
		g.GenerateNextRequest()
		if g.count > g.endPageNum {
			close(g.requestChan)
			return
		}
	}

}

func (g *Generator) GenerateStartRequest(seeds []string) {

	for _, url := range seeds {
		g.startRequests = append(g.startRequests, engine.Request{
			Url:        url,
			ParserFunc: parser.ParseTopicList,
			Agent:      uarand.GetRandom(),
		})
	}
	return
}

func (g *Generator) GenerateNextRequest() {
	var aux string
	if g.count == 0 {
		aux = ""
	} else {
		aux = nextString + strconv.Itoa(g.count)
	}
	for _, request := range g.startRequests {
		newRequest := request
		newRequest.Url = newRequest.Url + aux
		g.requestChan <- newRequest
	}
	g.count++

}

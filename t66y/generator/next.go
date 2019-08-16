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
	realms        []string
	startRequests []engine.Request
	count         int
	endPageNum    int
	requestChan   chan engine.Request
}

func NewGenerator(realms []string) chan engine.Request {
	start, stop := config.C.GetPageLimit()
	g := &Generator{
		realms:     realms,
		count:      start,
		endPageNum: stop,
	}
	g.requestChan = make(chan engine.Request)
	go g.Generate()
	return g.requestChan
}

func (g *Generator) SetRequestChan(requestChan chan engine.Request) {
	g.requestChan = requestChan

}
func (g *Generator) Generate() {
	g.GenerateStartRequest(g.realms)
	for {
		g.GenerateNextRequest()
		if g.count > g.endPageNum {
			close(g.requestChan)
			return
		}
	}

}

func (g *Generator) GenerateStartRequest(realms []string) {

	for _, url := range realms {
		g.startRequests = append(g.startRequests, engine.Request{
			Url:        url,
			ParserFunc: parser.ParsePostList,
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

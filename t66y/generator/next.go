package generator

import (
    "strconv"

    "github.com/corpix/uarand"

    "github.com/Tomilla/imagespider/config"
    "github.com/Tomilla/imagespider/engine"
    "github.com/Tomilla/imagespider/t66y/parser"
    "github.com/Tomilla/imagespider/util"
)

const nextString = "search=&page="

type Generator struct {
    realms        []string
    startRequests []engine.BaseParser
    count         int
    startPageNum  int
    stopPageNum   int
    requestChan   chan engine.BaseParser
}

func NewGenerator(realms []string) chan engine.BaseParser {
    start, stop := config.C.GetPageLimit()
    g := &Generator{
        realms:       realms,
        count:        0,
        startPageNum: start,
        stopPageNum:  stop,
    }
    g.requestChan = make(chan engine.BaseParser)
    go g.Generate()
    return g.requestChan
}

func (g *Generator) SetRequestChan(requestChan chan engine.BaseParser) {
    g.requestChan = requestChan

}
func (g *Generator) Generate() {
    g.GenerateStartRequest(g.realms)
    for {
        g.GenerateNextRequest()
        if g.count > g.stopPageNum {
            close(g.requestChan)
            return
        }
    }

}

func (g *Generator) GenerateStartRequest(realms []string) {

    for _, url := range realms {
        g.startRequests = append(g.startRequests, &parser.PostListRequest{
            URL:   url,
            Agent: uarand.GetRandom(),
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
        newRequest.SetURL(util.ConcatenateUrlOrder(newRequest.GetURL(), util.GetQueryPair(aux), []string{}))
        config.L.Infof("RequestsChan Url: %v", newRequest.GetURL())
        g.requestChan <- newRequest
    }
    g.count++
}

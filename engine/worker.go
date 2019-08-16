package engine

import (
	"log"

	"github.com/Tomilla/imagespider/fetcher"
)

type Worker struct {
}

func newWorker() *Worker {
	return &Worker{}
}

// fetch as request and return the parsed result

func (w *Worker) work(s Scheduler, out chan ParseResult) {
	const (
		logPrefix = "[engine worker] "
	)
	workChan := make(chan BaseParser)
	s.SubmitWorker(workChan)

	for {
		r := <-workChan
		url, post := r.GetURL(), r.GetPost()
		isList := post == nil

		if isList {
			log.Printf(logPrefix+"Fetching Post List, url: %s \n", url)
		} else {
			log.Printf(logPrefix+"Fetching Post %s, url: %s \n", post.Title, url)
		}

		body, err := fetcher.Fetch(url)
		if err != nil {
			if isList {
				log.Printf(logPrefix+"Fetching Post List error: %v %v\n", err, url)
			} else {
				log.Printf(logPrefix+"Fetching Post List error: %v %v %v\n", err, post, url)
			}
			s.SubmitWorker(workChan)
			continue
		}

		ParseResult := r.Parser(body, url)
		out <- ParseResult
		s.SubmitWorker(workChan)
	}

}

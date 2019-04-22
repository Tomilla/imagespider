package engine

import (
	"github.com/wuxiangzhou2010/imagespider/fetcher"
	"log"
)

type Worker struct {
}

func newWorker() *Worker {
	return &Worker{}
}

// fetch as request and return the parsed result

func (w *Worker) work(s Scheduler, out chan ParseResult) {
	workChan := make(chan Request)
	s.SubmitWorker(workChan)

	for {
		r := <-workChan

		log.Printf("Fetching %s, %s \n", r.Name, r.Url)
		body, err := fetcher.Fetch(r.Url)
		if err != nil {
			//panic(err)
			log.Println("Fetching error:", err, r.Name, r.Url)
			s.SubmitWorker(workChan)
			continue
		}

		ParseResult := r.ParserFunc(body, r.Url)
		out <- ParseResult
		s.SubmitWorker(workChan)
	}

}

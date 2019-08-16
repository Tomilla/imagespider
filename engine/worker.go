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
    workChan := make(chan Request)
    s.SubmitWorker(workChan)

    for {
        r := <-workChan

        log.Printf("[engine worker] Fetching %s, url: %s \n", r.Name, r.Url)
        body, err := fetcher.Fetch(r.Url)
        if err != nil {
            // panic(err)
            log.Println("[engine worker] Fetching error:", err, r.Name, r.Url)
            s.SubmitWorker(workChan)
            continue
        }

        ParseResult := r.ParserFunc(body, r.Url)
        out <- ParseResult
        s.SubmitWorker(workChan)
    }

}

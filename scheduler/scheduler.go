package scheduler

import "github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/engine"

type Scheduler struct {
	requestChan chan engine.Request
	workChan    chan chan engine.Request //可以接收Request
	WorkerCount int
}

func (s *Scheduler) Shutdown() {

}
func (s *Scheduler) GetWorkCount() int {
	return s.WorkerCount
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		requestChan: make(chan engine.Request),
		workChan:    make(chan chan engine.Request),
		WorkerCount: 10,
	}
}

func (s *Scheduler) Schedule() {

	var requestQ []engine.Request
	var workQ []chan engine.Request

	go func() {

		for {

			var activeWorker chan engine.Request
			var activeRequest engine.Request
			if len(requestQ) != 0 && len(workQ) != 0 {
				activeRequest = requestQ[0]
				activeWorker = workQ[0]
			}

			select {
			case newReq := <-s.requestChan:
				requestQ = append(requestQ, newReq)
			case readyWorker := <-s.workChan:
				workQ = append(workQ, readyWorker)
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workQ = workQ[1:]

			}
		}
	}()
}

func (s *Scheduler) SubmitRequest(request engine.Request) {
	s.requestChan <- request
}
func (s *Scheduler) SubmitWorker(worker chan engine.Request) {
	s.workChan <- worker
}

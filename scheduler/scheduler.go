package scheduler

import (
	"github.com/wuxiangzhou2010/imagespider/engine"
	"log"
)

type Scheduler struct {
	requestChan chan engine.Request
	workChan    chan chan engine.Request //可以接收Request

}

func (s *Scheduler) Shutdown() {

}

func NewScheduler() *Scheduler {
	return &Scheduler{
		requestChan: make(chan engine.Request),
		workChan:    make(chan chan engine.Request),
	}
}

func (s *Scheduler) Schedule(hungry chan bool) {

	var requestQ []engine.Request
	var workQ []chan engine.Request
	//tick := time.Tick(1 * time.Second)

	hungry <- true

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

				if len(requestQ) == 0 {
					hungry <- true
				}
			//case <-tick:
				log.Printf("[scheduler][requestQ len %d, cap %d], [workQ len %d, cap %d]\n",
					len(requestQ), cap(requestQ), len(workQ), cap(workQ))

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

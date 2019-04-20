package image

import "github.com/wuxiangzhou2010/imagespider/config"

type scheduler struct {
	workChan    chan work
	readyChan   chan chan work
	workerCount int
}

func newScheduler(workChan chan work, readyChan chan chan work) *scheduler {
	return &scheduler{workChan: workChan, readyChan: readyChan, workerCount: config.C.GetImageWorkerCount()}
}

func (s *scheduler) schedule() {
	var workQ []work
	var readyQ []chan work

	for {
		var activeWork work
		var activeWorker chan work
		if len(workQ) != 0 && len(readyQ) != 0 {
			activeWork = workQ[0]
			activeWorker = readyQ[0]
		}

		select {

		case w := <-s.workChan:
			workQ = append(workQ, w)
		case r := <-s.readyChan:
			readyQ = append(readyQ, r)
		case activeWorker <- activeWork:
			readyQ = readyQ[1:]
			workQ = workQ[1:]

		}
	}

}

func (s *scheduler) SubmitWork(w work) {
	s.workChan <- w
}
func (s *scheduler) Ready(c chan work) {
	s.readyChan <- c

}

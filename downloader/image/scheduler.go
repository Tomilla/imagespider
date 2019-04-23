package image

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/wuxiangzhou2010/imagespider/config"
)

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
	var preCount int32
	ticker := time.Tick(2 * time.Second)
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

		case <-ticker:
			v := atomic.LoadInt32(&count)
			if !atomic.CompareAndSwapInt32(&preCount, v, v) {
				preCount = v
				log.Printf("[Downloader worker] #%d downloaded [workQ len %d cap %d], [readyQ len %d cap %d]\n",
					v, len(workQ), cap(workQ), len(readyQ), cap(readyQ))
			}

			// 如果任务为空了， 则请求添加任务
			if len(workQ) == 0 && len(readyQ) == s.workerCount {
				ch := config.C.GetImageHungryChan()
				ch <- true
				fmt.Println("[All image requests are done, request more]")
				time.Sleep(5)
			}
		}
	}

}

func (s *scheduler) SubmitWork(w work) {
	s.workChan <- w
}
func (s *scheduler) Ready(c chan work) {
	s.readyChan <- c

}

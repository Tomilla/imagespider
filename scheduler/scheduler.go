package scheduler

import (
    "log"
    "sync/atomic"
    "time"

    "github.com/Tomilla/imagespider/config"
    "github.com/Tomilla/imagespider/engine"
)

type Scheduler struct {
    requestChan chan engine.BaseParser
    workChan    chan chan engine.BaseParser // 可以接收Request
}

func (s *Scheduler) Shutdown() {

}

func NewScheduler() *Scheduler {
    return &Scheduler{
        requestChan: make(chan engine.BaseParser),
        workChan:    make(chan chan engine.BaseParser),
    }
}

func (s *Scheduler) Schedule(hungry chan bool) {

    var requestQ []engine.BaseParser
    var workQ []chan engine.BaseParser
    tick := time.Tick(2 * time.Second)

    hungry <- true
    var count int32
    go func() {
        var preCount int32
        for {

            var activeWorker chan engine.BaseParser
            var activeRequest engine.BaseParser
            if len(requestQ) != 0 && len(workQ) != 0 {
                activeRequest = requestQ[0]
                activeWorker = workQ[0]
            }

            select {
            case newReq := <-s.requestChan:
                config.L.Infof("Url: %v", newReq.GetURL())
                requestQ = append(requestQ, newReq)
            case readyWorker := <-s.workChan:
                workQ = append(workQ, readyWorker)
            case activeWorker <- activeRequest:
                requestQ = requestQ[1:]
                workQ = workQ[1:]
            case <-tick:
                if len(requestQ) == 0 && len(workQ) == config.C.GetEngineWorkerCount() {
                    ch := config.C.GetImageHungryChan()
                    if ch == nil {
                        hungry <- true
                        time.Sleep(30 * time.Millisecond)
                    } else {
                        select {
                        case <-ch:
                            hungry <- true
                            time.Sleep(30 * time.Millisecond)
                        default:
                        }
                    }

                }
                v := atomic.LoadInt32(&count)
                if !atomic.CompareAndSwapInt32(&preCount, v, v) {
                    preCount = v
                    log.Printf("[scheduler][requestQ len %d, cap %d], [workQ len %d, cap %d]\n",
                        len(requestQ), cap(requestQ), len(workQ), cap(workQ))
                }
            }
        }
    }()
}

func (s *Scheduler) SubmitRequest(request engine.BaseParser) {
    s.requestChan <- request
}
func (s *Scheduler) SubmitWorker(worker chan engine.BaseParser) {
    s.workChan <- worker
}

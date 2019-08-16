package image

type Scheduler interface {
    schedule()
}
type Downloader interface {
    Run()
}

type Worker interface {
    work()
}

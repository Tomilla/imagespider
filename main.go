package main

import (
    "github.com/Tomilla/imagespider/t66y/generator"

    "github.com/Tomilla/imagespider/common"

    "github.com/Tomilla/imagespider/scheduler"

    "github.com/Tomilla/imagespider/engine"

    // for image init
    _ "github.com/Tomilla/imagespider/persist/image"
)

func main() {
    e := engine.NewConcurrentEngine(common.C.GetImageChan())
    e.Run(scheduler.NewScheduler(), generator.NewGenerator(common.C.GetStartPages()))

    // {
    //     osSignals := make(chan os.Signal, 1)
    //     signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
    //     <-osSignals
    // }
    {
        e.Shutdown()
    }
}

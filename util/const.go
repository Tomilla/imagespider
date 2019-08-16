package util

import (
    "log"
    "time"
)

const (
    DateLayout      = "2006-01-02"
    DateDefault     = "3000-01-01"
    DateTimeLayout  = "2006-01-02 15:04"
    DateTimeDefault = "3000-01-01 01:01:00"
    TimeZone        = "Asia/Shanghai"
    DefaultFilePerm = 0755
)

var (
    TZ *time.Location
)

func init() {
    var err error
    TZ, err = time.LoadLocation(TimeZone)
    if err != nil {
        log.Panic("cannot load time-zone")
    }
}

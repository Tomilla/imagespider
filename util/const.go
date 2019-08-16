package util

import (
    "log"
    "time"
)

var DateLayout = "2006-01-02"
var DateDefault = "3000-01-01"
var DateTimeLayout = "2006-01-02 15:04"
var DateTimeDefault = "3000-01-01 01:01:00"
var TimeZone = "Asia/Shanghai"
var TZ *time.Location

func init() {
    var err error
    TZ, err = time.LoadLocation(TimeZone)
    if err != nil {
        log.Panic("cannot load time-zone")
    }
}

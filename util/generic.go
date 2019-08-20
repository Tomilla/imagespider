package util

import (
    "fmt"
    "reflect"

    "github.com/Tomilla/imagespider/common"
)

func EqualSliceGeneric(lhs interface{}, rhs interface{}) (isEqual bool) {
    tl, tr := reflect.TypeOf(lhs), reflect.TypeOf(rhs)
    vl, vr := reflect.ValueOf(lhs), reflect.ValueOf(rhs)
    if tl.Kind() != reflect.Slice || tl.Kind() != reflect.Slice || tl != tr {
        return false
    }
    if vl.Len() != vr.Len() {
        return false
    }
    ch := make(chan bool, 1)
    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered in f", r)
                isEqual = false
            }
            close(ch)
        }()
        for i := 0; i < vl.Len(); i++ {
            el, er := vl.Index(i), vr.Index(i)
            switch el.Kind() {
            case reflect.Bool:
                if el.Bool() != er.Bool() {
                    ch <- false
                    return
                }
            case reflect.String:
                if el.String() != er.String() {
                    ch <- false
                    return
                }
            default:
                ch <- false
                return
            }
        }
        ch <- true
    }()
    isEqual = <-ch
    common.L.Infof("isEqual: %v", isEqual)
    return
}

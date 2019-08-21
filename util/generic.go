package util

import (
    "fmt"
    "reflect"
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
                fmt.Println("Recovered in go-routine: ", r)
                isEqual = false
                ch <- isEqual
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
            case reflect.Int:
                fallthrough
            case reflect.Int8:
                fallthrough
            case reflect.Int16:
                fallthrough
            case reflect.Int32:
                fallthrough
            case reflect.Int64:
                // Int returns v's underlying value, as a largest possible int variant -- int64.
                if el.Int() != er.Int() {
                    ch <- false
                    return
                }
                break
            case reflect.Uint:
                fallthrough
            case reflect.Uint8:
                fallthrough
            case reflect.Uint16:
                fallthrough
            case reflect.Uint32:
                fallthrough
            case reflect.Uint64:
                // Uint returns v's underlying value, as a largest possible uint variant --  uint64.
                if el.Uint() != er.Uint() {
                    ch <- false
                    return
                }
                break
            case reflect.Uintptr:
                if el.Pointer() != er.Pointer() {
                    ch <- false
                    return
                }
                break
            case reflect.Float32:
                fallthrough
            case reflect.Float64:
                // Float returns v's underlying value, as a largest possible float variant  float64e reflejct.Complex64:
                if el.Float() != er.Float() {
                    ch <- false
                    return
                }
                break
            case reflect.Complex64:
                fallthrough
            case reflect.Complex128:
                if el.Complex() != er.Complex() {
                    ch <- false
                    return
                }
                break
            default:
                ch <- false
                return
            }
        }
        ch <- true
    }()
    isEqual = <-ch
    // common.L.Infof("isEqual: %v", isEqual)
    return
}

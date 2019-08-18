package util

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)

func CheckErr(e error) {
    if e != nil {
        panic(e)
    }
}

func SleepRandomDuration(min int, max int) {
    t := rand.Intn(max-min+1) + min
    // fmt.Printf("Sleep: %v\n", t)
    time.Sleep(time.Duration(t) * time.Millisecond)
}

func WarnErr(e error) {
    if e != nil {
        fmt.Printf("Warn: %v\n", e)
    }
}

func GetLastItem(arr interface{}) string {
    switch arr.(type) {
    case [][]string:
        _arr := arr.([][]string)
        outerLength := len(_arr)
        if outerLength > 0 {
            return GetLastItem(_arr[outerLength-1])
        }
        break
    case []string:
        _arr := arr.([]string)
        innerLength := len(_arr)
        if innerLength > 0 {
            return _arr[innerLength-1]
        }
        break
    default:
        return ""
    }
    return ""
}

func ConcatenateUrl(url string, iterator map[string]string, exclude []string) string {
    // sort.Slice(exclude, func (i, j int) bool {
    //     return exclude[i] < exclude[j]
    // })
    // needle := "test"
    // idx := sort.Search(len(exclude), func (i int) bool {
    //     return exclude[i] == needle
    // })
    var (
        existentQuery = map[string]string{}
        junction      = "?"
        requestArgs   []string
        pushFunc      = func(option map[string]string) {
            key := option["key"]
            val := option["val"]
            if !Contains(exclude, key) {
                requestArgs = append(requestArgs, key+"="+val)
            }
        }
        pushFuncWrapper = func(reversed bool) func(string, string) {
            return func(arg1, arg2 string) {
                if reversed {
                    arg1, arg2 = arg2, arg1
                }
                pushFunc(map[string]string{
                    "key": arg1,
                    "val": arg2,
                })
            }
        }
    )

    if strings.Contains(url, "?") {
        urlParts := strings.SplitN(url, "?", 2)
        url = urlParts[0]
        query := urlParts[1]

        if len(query) > 0 {
            queryPairs := strings.SplitN(query, "&", -1)
            for _, pair := range queryPairs {
                _pair := strings.SplitN(pair, "=", 2)
                existentQuery[_pair[0]] = _pair[1]
            }
        }
    }

    f := pushFuncWrapper(false)
    for key, val := range iterator {
        _, ok := existentQuery[key]
        if ok {
            fmt.Printf("Override: %v -> %v", existentQuery[key], val)
            existentQuery[key] = val
        } else {
            existentQuery[key] = val
        }
    }
    for key, val := range existentQuery {
        f(key, val)
    }
    return url + junction + strings.Join(requestArgs, "&")
}

package util

import (
    "fmt"
    "math/rand"
    "strings"
    "time"

    "github.com/Tomilla/imagespider/collections/set"
    "github.com/Tomilla/imagespider/config"
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

func ConcatenateUrlOrder(url string, iterator [][]string, exclude []string) string {
    var (
        existentQuery [][]string
        junction            = "?"
        requestArgs   []string
        pushFunc            = func(option map[string]string) {
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
                existentQuery = append(existentQuery, strings.SplitN(pair, "=", 2))
            }
        }
    }

    f := pushFuncWrapper(false)
    queryMap := NewOrderQueryMap()
    queryMap.Init(existentQuery)
    for _, item := range iterator {
        key, val := item[0], item[1]
        _, ok := queryMap.M[key]
        if ok {
            config.L.Debugf("Override: %v -> %v", queryMap.M[key], val)
            queryMap.Set(key, val)
        } else {
            queryMap.Set(key, val)
        }
    }
    queryMap.Iterate(f)
    return url + junction + strings.Join(requestArgs, "&")
}

func ConcatenateUrl(url string, iterator map[string]string, exclude []string) string {
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
            config.L.Infof("Override: %v -> %v", existentQuery[key], val)
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

func GetQuerySet(url string) *set.Set {
    querySet := set.New()
    if strings.Contains(url, queryStartMark) {
        queryPart := strings.SplitN(url, queryStartMark, 2)[1]
        for _, pair := range strings.SplitN(queryPart, querySeparator, -1) {
            querySet.Insert(pair)
        }
    }
    return querySet
}

func GetQueryPair(url string) [][]string {
    var (
        queryPart string
        result    [][]string
    )
    if len(url) == 0 {
        return [][]string{}
    }

    if strings.Contains(url, queryStartMark) {
        queryPart = strings.SplitN(url, queryStartMark, 2)[1]
    } else {
        queryPart = url
    }
    for _, pair := range strings.SplitN(queryPart, querySeparator, -1) {
        result = append(result, strings.SplitN(pair, queryPairMark, 2))
    }
    return result
}

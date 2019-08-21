package util

import (
    "regexp"
    "strings"
)

var (
    WhiteSpace = " \t\n\r\x0b\x0c"
    Contains   = StringInSlice
    Find       = StringInSlice
)

// https://github.com/DaddyOh/golang-samples/blob/master/pad.go
// RightPad2Len
func RightPad2Len(s string, padStr string, width int) string {
    var padCountInt = 1 + ((width - len(padStr)) / len(padStr))
    var retStr = s + strings.Repeat(padStr, padCountInt)
    return retStr[:width]
}

// LeftPad2Len
func LeftPad2Len(s string, padStr string, width int) string {
    var padCountInt = 1 + ((width - len(padStr)) / len(padStr))
    var retStr = strings.Repeat(padStr, padCountInt) + s
    return retStr[(len(retStr) - width):]
}

// StringInSlice
func StringInSlice(list []string, a string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func StringIndexSlice(a []string, x string) int {
    for i, n := range a {
        if x == n {
            return i
        }
    }
    return len(a)
}

func getCorrectRegexObj(regEx interface{}) (compRegEx *regexp.Regexp) {
    switch regEx.(type) {
    case string:
        regEx := regEx.(string)
        compRegEx = regexp.MustCompile(regEx)
    case *regexp.Regexp:
        regEx := regEx.(*regexp.Regexp)
        compRegEx = regEx
    default:
    }
    return
}

/**
 * Parses string with the given regular expression and returns the
 * group values defined in the expression.
 */
func GetRegexNamedGroupMapping(regEx interface{}, txt string) (paramsMap map[string]string) {
    var (
        compRegEx *regexp.Regexp
    )
    compRegEx = getCorrectRegexObj(regEx)
    if compRegEx == nil {
        return
    }

    paramsMap = make(map[string]string)
    match := compRegEx.FindStringSubmatch(txt)
    for i, name := range compRegEx.SubexpNames() {
        if i > 0 && i <= len(match) {
            paramsMap[name] = match[i]
        }
    }
    return
}

/**
 * Same as above but normal number-based group
 */
func GetRegexGroupArray(regEx interface{}, txt string) (paramsArray []string) {
    var (
        compRegEx *regexp.Regexp
    )
    compRegEx = getCorrectRegexObj(regEx)
    if compRegEx == nil {
        return
    }

    paramsArray = compRegEx.FindStringSubmatch(txt)
    if len(paramsArray) > 1 {
        return paramsArray[1:]
    }
    return []string{}
}

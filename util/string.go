package util

import (
	"strings"
)

var WhiteSpace string = " \t\n\r\x0b\x0c"

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
// func StringInSlice(list []string, a string) bool {
//     for _, b := range list {
//         if b == a {
//             return true
//         }
//     }
//     return false
// }

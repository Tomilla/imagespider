package util

import (
	"fmt"
	"math/rand"
	"time"
)

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func SleepRandomDuration(min int, max int) {
	t := rand.Intn(max-min+1) + min
	// fmt.Printf("Seleep: %v\n", t)
	time.Sleep(time.Duration(t))
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


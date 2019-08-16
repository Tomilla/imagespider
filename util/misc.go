package util

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func SleepRandomDuration(min int, max int) {
	t := rand.Intn(max-min+1) + min
	fmt.Printf("Seleep: %v\n", t)
	time.Sleep(time.Duration(t))
}

func WarnErr(e error) {
	if e != nil {
		fmt.Printf("Warn: %v\n", e)
	}
}

func CheckPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

package util

import (
	"os"
)

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

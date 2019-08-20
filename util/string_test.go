package util

import (
    "regexp"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGetRegexGroupArray(t *testing.T) {
    type TestInfo struct {
        name     string
        regex    interface{}
        excepted []string
    }

    var (
        sourceAndExpected = []TestInfo{
            // test with normal string
            {"3604380.html", `(\d+)(\.\w+)`, []string{"3604380", ".html"}},
            {"47920dx8.txt", `(\w+)(\.\w+)`, []string{"47920dx8", ".txt"}},
            // test with regex obj
            {"resource.php", regexp.MustCompile(`(?i)([a-z]+)\.([a-z]+)`), []string{"resource", "php"}},
        }
    )
    for _, info := range sourceAndExpected {
        path, regex, excepted := info.name, info.regex, info.excepted
        result := GetRegexGroupArray(regex, path)
        t.Logf("%v", result)
        t.Logf("%v", excepted)
        assert.Equal(t, EqualSliceGeneric(result, excepted), true)
    }
}

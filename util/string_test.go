package util

import (
    "reflect"
    "regexp"
    "testing"

    "github.com/google/go-cmp/cmp"
    "github.com/stretchr/testify/assert"
)

func TestGetRegexGroupArray(t *testing.T) {
    type TestInfo struct {
        name     string
        regex    interface{}
        excepted []string
    }

    var (
        pureLetterRe      = regexp.MustCompile(`(?i)([a-z]+)\.([a-z]+)`)
        sourceAndExpected = []TestInfo{
            // test with normal string
            {"3604380.html", `(\d+)(\.\w+)`, []string{"3604380", ".html"}},
            {"47920dx8.txt", `(\w+)(\.\w+)`, []string{"47920dx8", ".txt"}},
            // test with regex obj
            {"resource.php", pureLetterRe, []string{"resource", "php"}},
            {"resource.PHP", pureLetterRe, []string{"resource", "PHP"}},
            {"47920dx8.txt", regexp.MustCompile(`(?i)(\d{2,})(\.\w+)`), []string{}},
            {"47920dx8.txt", regexp.MustCompile(`(?i)(\d{2,})(\.\w+)`), nil},
        }
    )
    for _, info := range sourceAndExpected {
        path, regex, excepted := info.name, info.regex, info.excepted
        result := GetRegexGroupArray(regex, path)
        t.Logf("actual: %v", result)
        t.Logf("excepted: %v", excepted)
        assert.Equal(t, true, EqualSliceGeneric(result, excepted))
        if excepted != nil {
            assert.Equal(t, true, reflect.DeepEqual(result, excepted))
        } else {
            assert.NotEqual(t, true, reflect.DeepEqual(result, excepted))
        }
    }
}

func TestGetRegexNamedGroupMapping(t *testing.T) {
    type TestInfo struct {
        name     string
        regex    interface{}
        excepted map[string]string
    }

    var (
        pureLetterRe      = regexp.MustCompile(`(?i)(?P<Name>[a-z]+)\.(?P<Ext>[a-z]+)`)
        sourceAndExpected = []TestInfo{
            // test with normal string
            {"3604380.html", `(?P<Name>\d+)(?P<Ext>\.\w+)`, map[string]string{"Name": "3604380", "Ext": ".html"}},
            {"47920dx8.txt", `(?P<Name>\w+)(?P<Ext>\.\w+)`, map[string]string{"Name": "47920dx8", "Ext": ".txt"}},
            // test with regex obj
            {"resource.php", pureLetterRe, map[string]string{"Name": "resource", "Ext": "php"}},
            {"resource.PHP", pureLetterRe, map[string]string{"Name": "resource", "Ext": "PHP"}},
            {"47920dx8.txt", regexp.MustCompile(`(?i)(\d{2,})(\.\w+)`), map[string]string{}},
        }
    )

    for _, info := range sourceAndExpected {
        path, regex, excepted := info.name, info.regex, info.excepted
        result := GetRegexNamedGroupMapping(regex, path)
        t.Logf("actual: %v", result)
        t.Logf("excepted: %v", excepted)
        assert.Equal(t, cmp.Equal(result, excepted), true)
        assert.Equal(t, reflect.DeepEqual(result, excepted), true)
    }
}

func TestLeftPad2Len(t *testing.T) {
    type TestInfo struct {
        s        string
        p        string
        w        int
        excepted string
    }
    var (
        sourceAndExcepted = []TestInfo{
            {"1", "0", 10, "0000000001"},
            {"5", "0", 5, "00005"},
            {"5", " ", 5, "    5"},
            // empty pad string
            {"05", "", 5, "05"},
            // illegal value: width, return original string
            {"005", " ", -1, "005"},
            {"005", " ", 0, "005"},
        }
    )
    for _, info := range sourceAndExcepted {
        assert.Equal(t, LeftPad2Len(info.s, info.p, info.w), info.excepted)
    }
}

func TestRightPad2Len(t *testing.T) {
    type TestInfo struct {
        s        string
        p        string
        w        int
        excepted string
    }
    var (
        sourceAndExcepted = []TestInfo{
            {"1", "0", 10, "1000000000"},
            {"5", "0", 5, "50000"},
            {"5.", "3", 5, "5.333"},
            {"5", " ", 5, "5    "},
            // empty pad string
            {"05", "", 5, "05"},
            // illegal value: width, return original string
            {"005", " ", -1, "005"},
            {"005", " ", 0, "005"},
        }
    )
    for _, info := range sourceAndExcepted {
        assert.Equal(t, RightPad2Len(info.s, info.p, info.w), info.excepted)
    }

}
package util

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestConcatenateUrl(t *testing.T) {
    type TestInfo struct {
        url      string
        query    map[string]string
        exclude  []string
        excepted string
    }
    var (
        sourceAndExcepted = []TestInfo{
            {"http://t66y.com/thread0806.php?fid=16",
                map[string]string{"search": "", "page": "1"},
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=&page=1",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=1",
                map[string]string{"search": "", "page": "2"},
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=&page=2",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=2",
                map[string]string{"search": "A", "page": "3"},
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=A&page=3",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=3",
                map[string]string{"search": "B", "page": "4"},
                []string{"search"},
                "http://t66y.com/thread0806.php?fid=16&page=4",
            },
        }
    )
    for _, info := range sourceAndExcepted {
        url, query, exclude, excepted := info.url, info.query, info.exclude, info.excepted
        result := ConcatenateUrl(url, query, exclude)
        println(excepted)
        println(result)
        assert.Equal(t, excepted, result)
    }
}

package util

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestConcatenateUrlOrder(t *testing.T) {
    type TestInfo struct {
        url      string
        query    [][]string
        exclude  []string
        excepted string
    }
    var (
        sourceAndExcepted = []TestInfo{
            {"http://t66y.com/thread0806.php?fid=16",
                [][]string{
                    {"search", "",},
                    {"page", "1",},
                },
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=&page=1",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=1",
                [][]string{
                    {"search", "",},
                    {"page", "2",},
                },
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=&page=2",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=2",
                [][]string{
                    {"search", "A",},
                    {"page", "3",},
                },
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=A&page=3",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=3",
                [][]string{
                    {"search", "B",},
                    {"page", "4",},
                },
                []string{"search"},
                "http://t66y.com/thread0806.php?fid=16&page=4",
            },
        }
    )
    for _, info := range sourceAndExcepted {
        url, query, exclude, excepted := info.url, info.query, info.exclude, info.excepted
        result := ConcatenateUrlOrder(url, query, exclude)
        println(excepted)
        println(result)
        assert.Equal(t, excepted, result)
    }
}

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
        sourceAndWrong = []TestInfo{
            {"http://t66y.com/thread0806.php?fid=16",
                map[string]string{"search": "", "page": "1"},
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=&page=1&hello=tomi",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=3",
                map[string]string{"search": "B", "page": "4"},
                []string{"search"},
                "http://t66y.com/thread0806.php?fid=18&page=3",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=1",
                map[string]string{"search": "", "page": "2"},
                []string{},
                "http://t66y.com/thread0806.php?fid=15",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=2",
                map[string]string{"search": "B", "page": "4"},
                []string{},
                "http://t66y.com/thread0806.php?fid=16&search=&page=2",
            },
            {"http://t66y.com/thread0806.php?fid=16&search=&page=2",
                map[string]string{"search": "B", "page": "4"},
                []string{},
                "http://t66y.com/thread0806.php",
            },
        }
        startTest = func(tInfo []TestInfo, isEqual bool) {
            for _, info := range tInfo {
                url, query, exclude, excepted := info.url, info.query, info.exclude, info.excepted
                result := ConcatenateUrl(url, query, exclude)
                set1 := GetQuerySet(result)
                set2 := GetQuerySet(excepted)
                if isEqual {
                    assert.Equal(t, set1.Len(), set2.Len())
                    assert.Equal(t, set1.Intersection(set2).Len(), set2.Len())
                } else {
                    println(result)
                    println(excepted)
                    if set2.Len() > 0 {
                        assert.NotEqual(t, set2.Len(), set1.Intersection(set2).Len())
                    } else {
                        assert.NotEqual(t, 0, set1.Len())
                    }
                }
            }
        }
    )
    startTest(sourceAndExcepted, true)
    startTest(sourceAndWrong, false)
}

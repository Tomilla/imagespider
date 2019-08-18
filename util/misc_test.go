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
            {"http://example.com/thread0806.php?fid=16",
                [][]string{
                    {"search", "",},
                    {"page", "1",},
                },
                []string{},
                "http://example.com/thread0806.php?fid=16&search=&page=1",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=1",
                [][]string{
                    {"search", "",},
                    {"page", "2",},
                },
                []string{},
                "http://example.com/thread0806.php?fid=16&search=&page=2",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=2",
                [][]string{
                    {"search", "A",},
                    {"page", "3",},
                },
                []string{},
                "http://example.com/thread0806.php?fid=16&search=A&page=3",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=3",
                [][]string{
                    {"search", "B",},
                    {"page", "4",},
                },
                []string{"search"},
                "http://example.com/thread0806.php?fid=16&page=4",
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
            {"http://example.com/thread0806.php?fid=16",
                map[string]string{"search": "", "page": "1"},
                []string{},
                "http://example.com/thread0806.php?fid=16&search=&page=1",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=1",
                map[string]string{"search": "", "page": "2"},
                []string{},
                "http://example.com/thread0806.php?fid=16&search=&page=2",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=2",
                map[string]string{"search": "A", "page": "3"},
                []string{},
                "http://example.com/thread0806.php?fid=16&search=A&page=3",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=3",
                map[string]string{"search": "B", "page": "4"},
                []string{"search"},
                "http://example.com/thread0806.php?fid=16&page=4",
            },
        }
        sourceAndWrong = []TestInfo{
            {"http://example.com/thread0806.php?fid=16",
                map[string]string{"search": "", "page": "1"},
                []string{},
                "http://example.com/thread0806.php?fid=16&search=&page=1&hello=tomi",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=3",
                map[string]string{"search": "B", "page": "4"},
                []string{"search"},
                "http://example.com/thread0806.php?fid=18&page=3",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=1",
                map[string]string{"search": "", "page": "2"},
                []string{},
                "http://example.com/thread0806.php?fid=15",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=2",
                map[string]string{"search": "B", "page": "4"},
                []string{},
                "http://example.com/thread0806.php?fid=16&search=&page=2",
            },
            {"http://example.com/thread0806.php?fid=16&search=&page=2",
                map[string]string{"search": "B", "page": "4"},
                []string{},
                "http://example.com/thread0806.php",
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

func TestGetQuerySet(t *testing.T) {
    type TestInfo struct {
        url  string
        size int
    }
    var (
        sourceAndExcepted = []*TestInfo{
            {"http://example.com/thread0806.php?fid=16", 1},
            {"http://example.com/thread0806.php?fid=16&search=", 2},
            {"http://example.com/thread0806.php?fid=16&search=&page=1", 3},
            {"http://example.com/thread0806.php?fid=18&search=B&page=2", 3},
            {"http://example.com/thread0806.php?fid=18&search=B&page=2&name=tomi", 4},
        }
        sourceAndWrong = []*TestInfo{
            {"http://example.com/thread0806.php?fid=16", 2},
            {"http://example.com/thread0806.php?fid=16&search=", 4},
            {"http://example.com/thread0806.php?fid=16&search=&page=1", 1},
            {"http://example.com/thread0806.php?fid=18&search=B&page=2", 2},
            {"http://example.com/thread0806.php?fid=18&search=B&page=2&name=tomi", 5},
        }
    )
    for _, info := range sourceAndExcepted {
        assert.Equal(t, GetQuerySet(info.url).Len(), info.size)
    }
    for _, info := range sourceAndWrong {
        assert.NotEqual(t, GetQuerySet(info.url).Len(), info.size)
    }
}

func TestGetQueryPair(t *testing.T) {
    type TestInfo struct {
        url string
        ret [][]string
    }
    var (
        sourceAndExcepted = []*TestInfo{
            {"http://example.com/thread0806.php?fid=16", [][]string{
                {"fid", "16",},
            }},
            {"http://example.com/thread0806.php?fid=16&search=", [][]string{
                {"fid", "16",},
                {"search", "",},
            }},
            {"http://example.com/thread0806.php?fid=16&search=&page=1", [][]string{
                {"fid", "16",},
                {"search", "",},
                {"page", "1",},
            }},
            {"http://example.com/thread0806.php?fid=18&search=B&page=2", [][]string{
                {"fid", "18",},
                {"search", "B",},
                {"page", "2",},
            }},
            {"http://example.com/thread0806.php?fid=18&search=B&page=2&name=tomi", [][]string{
                {"fid", "18",},
                {"search", "B",},
                {"page", "2",},
                {"name", "tomi",},
            }},
            // test without host part and path part
            {"fid=18&search=B&page=2", [][]string{
                {"fid", "18",},
                {"search", "B",},
                {"page", "2",},
            }},
            // test empty string
            {"", [][]string{
            }},
        }
        sourceAndWrong = []*TestInfo{
            // test different pair
            {"http://example.com/thread0806.php?fid=16", [][]string{
                {"fid", "18",},
            }},
            // test different length
            {"http://example.com/thread0806.php?fid=16&search=", [][]string{
                {"fid", "16",},
            }},
            // same as above, but except longer length
            {"http://example.com/thread0806.php?fid=16&search=&page=1", [][]string{
                {"fid", "16",},
                {"search", "",},
                {"page", "1",},
                {"limit", "1",},
            }},
            // test wrong order
            {"http://example.com/thread0806.php?fid=18&search=B&page=2", [][]string{
                {"fid", "12",},
                {"page", "2",},
                {"search", "A",},
            }},
            // test wrong size and wrong value
            {"http://example.com/thread0806.php?fid=18&search=B&page=2&name=tomi", [][]string{
                {"fid", "18",},
                {"search", "B",},
                {"page", "1",},
            }},
            // test without host part and path part
            {"fid=18&search=B&page=2", [][]string{
                {"fid", "18",},
                {"search", "B",},
                {"page", "3",},
            }},
            // test empty string
            {"", [][]string{
                {"search", "B",},
                {"page", "3",},
            }},
        }
    )
    for _, info := range sourceAndExcepted {
        result := GetQueryPair(info.url)
        assert.Equal(t, info.ret, result)
    }
    for _, info := range sourceAndWrong {
        result := GetQueryPair(info.url)
        assert.NotEqual(t, info.ret, result)
    }
}

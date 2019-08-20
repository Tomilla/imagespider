package util

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEqualSliceGeneric(t *testing.T) {
    type TestInfo struct {
        lhs      interface{}
        rhs      interface{}
        excepted bool
    }
    var (
        sourceAndExcepted = []TestInfo{
            // test: string slice
            {[]string{"hello", "world"}, []string{"hello", "world"}, true,},
            {[]string{"hello", "world"}, []string{"hello", "tomi"}, false,},
            // test: boolean slice
            {[]bool{true, true}, []bool{true, true}, true,},
            {[]bool{false, false}, []bool{false, true}, false,},
        }
    )
    for _, info := range sourceAndExcepted {
        assert.Equal(t, EqualSliceGeneric(info.lhs, info.rhs), info.excepted)
    }
}

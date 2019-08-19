package generator

import (
    "testing"
    "time"

    "github.com/Tomilla/imagespider/common"
)

func TestGenerator_GenerateStartRequest(t *testing.T) {
    ch := NewGenerator(common.C.GetStartPages())
    go func() {
        for v := range ch {
            t.Logf("%+v", v)
        }
    }()
    time.Sleep(100 * time.Millisecond)
}

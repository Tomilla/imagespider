package generator

import (
	"github.com/wuxiangzhou2010/imagespider/config"
	"testing"
	"time"
)

func TestGenerator_GenerateStartRequest(t *testing.T) {
	ch := NewGenerator(config.C.GetStartPages())
	go func() {
		for i := 0; i < 10; i++ {
			t.Logf("%+v", <-ch)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

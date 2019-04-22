package generator

import (
	"github.com/wuxiangzhou2010/imagespider/config"
	"testing"
	"time"
)

func TestGenerator_GenerateStartRequest(t *testing.T) {
	ch := NewGenerator(config.C.GetStartPages())
	go func() {
		for v := range ch {
			t.Logf("%+v", v)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

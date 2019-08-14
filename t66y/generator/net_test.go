package generator

import (
	"testing"
	"time"

	"github.com/wuxiangzhou2010/imagespider/config"
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

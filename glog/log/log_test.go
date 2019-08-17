package log

import (
    "bytes"
    "os"
    "path/filepath"
    "testing"

    "github.com/Tomilla/imagespider/glog"
)

type memoryWriter struct {
    bytes.Buffer
}

func (w *memoryWriter) Close() error {
    w.Buffer.Reset()
    return nil
}

func TestLogFuncs(t *testing.T) {
    w := &memoryWriter{}
    h := glog.NewHandler(glog.InfoLevel, glog.DefaultFormatter)
    h.AddWriter(w)
    l := glog.NewLogger(glog.InfoLevel)
    l.AddHandler(h)
    SetDefaultLogger(l)

    l.Debug("test")
    if w.Buffer.Len() != 0 {
        t.Error("memoryWriter is not empty")
    }
    l.Debugf("test")
    if w.Buffer.Len() != 0 {
        t.Error("memoryWriter is not empty")
    }

    l.Info("test")
    if w.Buffer.Len() == 0 {
        t.Error("memoryWriter is empty")
    }
    w.Buffer.Reset()

    l.Infof("test")
    if w.Buffer.Len() == 0 {
        t.Error("memoryWriter is empty")
    }
    w.Buffer.Reset()

    l.Error("test")
    if w.Buffer.Len() == 0 {
        t.Error("memoryWriter is empty")
    }
    w.Buffer.Reset()

    l.Errorf("test")
    if w.Buffer.Len() == 0 {
        t.Error("memoryWriter is empty")
    }
    l.Close()

    h = glog.NewHandler(glog.ErrorLevel, glog.DefaultFormatter)
    h.AddWriter(w)
    l = glog.NewLogger(glog.ErrorLevel)
    l.AddHandler(h)
    SetDefaultLogger(l)

    l.Info("test")
    if w.Buffer.Len() != 0 {
        t.Error("memoryWriter is not empty")
    }
    w.Buffer.Reset()

    l.Error("test")
    if w.Buffer.Len() == 0 {
        t.Error("memoryWriter is empty")
    }
    l.Close()
}

func BenchmarkBufferedFileLogger(b *testing.B) {
    path := filepath.Join(os.TempDir(), "test.log")
    err := os.Remove(path)
    if err != nil {
        b.Logf("cannot remove path: %v", path)
    }
    w, err := glog.NewBufferedFileWriter(path)
    if err != nil {
        b.Error(err)
    }
    h := glog.NewHandler(glog.InfoLevel, glog.DefaultFormatter)
    h.AddWriter(w)
    l := glog.NewLogger(glog.InfoLevel)
    l.AddHandler(h)
    SetDefaultLogger(l)

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Infof("test")
        }
    })
    l.Close()
}

func BenchmarkDiscardLogger(b *testing.B) {
    w := glog.NewDiscardWriter()
    h := glog.NewHandler(glog.InfoLevel, glog.DefaultFormatter)
    h.AddWriter(w)
    l := glog.NewLogger(glog.InfoLevel)
    l.AddHandler(h)
    SetDefaultLogger(l)

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Infof("test")
        }
    })
    l.Close()
}

func BenchmarkNopLog(b *testing.B) {
    w := glog.NewDiscardWriter()
    h := glog.NewHandler(glog.InfoLevel, glog.DefaultFormatter)
    h.AddWriter(w)
    l := glog.NewLogger(glog.InfoLevel)
    l.AddHandler(h)
    SetDefaultLogger(l)

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Debugf("test")
        }
    })
    l.Close()
}

func BenchmarkMultiLevels(b *testing.B) {
    w := glog.NewDiscardWriter()
    dh := glog.NewHandler(glog.DebugLevel, glog.DefaultFormatter)
    dh.AddWriter(w)
    ih := glog.NewHandler(glog.InfoLevel, glog.DefaultFormatter)
    ih.AddWriter(w)
    wh := glog.NewHandler(glog.WarnLevel, glog.DefaultFormatter)
    wh.AddWriter(w)
    eh := glog.NewHandler(glog.ErrorLevel, glog.DefaultFormatter)
    eh.AddWriter(w)
    ch := glog.NewHandler(glog.CritLevel, glog.DefaultFormatter)
    ch.AddWriter(w)

    l := glog.NewLogger(glog.WarnLevel)
    l.AddHandler(dh)
    l.AddHandler(ih)
    l.AddHandler(wh)
    l.AddHandler(eh)
    l.AddHandler(ch)
    SetDefaultLogger(l)

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Debugf("test")
            Infof("test")
            Warnf("test")
            Errorf("test")
            Critf("test")
        }
    })
    l.Close()
}

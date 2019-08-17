package glog

import (
    "errors"
    "io"
    "os"
    "path/filepath"
    "sync"
    "testing"
    "time"
    "unicode/utf8"

    "github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
    SetInternalLogger(NewStderrLogger())
    os.Exit(m.Run())
}

func TestBufferedFileWriter(t *testing.T) {
    var err error
    oldBufferSize := bufferSize
    bufferSize = 1024
    path := filepath.Join(os.TempDir(), "test.log")
    err = os.Remove(path)
    if err != nil {
        t.Errorf("Cannot remove path: %v\n", path)
    }
    w, err := NewBufferedFileWriter(path)
    if err != nil {
        t.Error(err)
    }

    f, err := os.Open(path)
    if err != nil {
        t.Error(err)
    }
    stat, err := f.Stat()
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    n, err := w.Write([]byte("test"))
    if err != nil {
        t.Error(err)
    }
    if n != 4 {
        t.Errorf("read %d bytes", n)
    }

    buf := make([]byte, bufferSize*2)
    n, err = f.Read(buf)
    if err != io.EOF {
        t.Error(err)
    }
    if n != 0 {
        t.Errorf("read %d bytes", n)
    }

    time.Sleep(flushDuration * 2)
    n, err = f.Read(buf)
    if err != nil {
        t.Error(err)
    }
    if n != 4 {
        t.Errorf("read %d bytes", n)
    }
    bs := string(buf[:4])
    if bs != "test" {
        t.Error("read bytes are " + bs)
    }

    for i := 0; i < bufferSize; i++ {
        _, err := w.Write([]byte{'1'})
        if err != nil {
            t.Error("Cannot write into w\n")
        }
    }
    _, err = w.Write([]byte{'2'}) // writes over bufferSize cause flushing
    if err != nil {
        t.Error("Cannot write into w\n")
    }
    n, err = f.Read(buf)
    if err != nil {
        t.Error(err)
    }
    if n != bufferSize {
        t.Errorf("read %d bytes", n)
    }
    if buf[bufferSize-1] != '1' {
        t.Errorf("last byte is %d", buf[bufferSize-1])
    }
    if buf[bufferSize] != 0 {
        t.Errorf("next byte is %d", buf[bufferSize-1])
    }

    time.Sleep(flushDuration * 2)
    n, err = f.Read(buf)
    if err != nil {
        t.Error(err)
    }
    if n != 1 {
        t.Errorf("read %d bytes", n)
    }
    if buf[0] != '2' {
        t.Errorf("first byte is %d", buf[0])
    }
    if buf[1] != '1' {
        t.Errorf("next byte is %d", buf[1])
    }

    err = f.Close()
    if err != nil {
        t.Error("Cannot close f")
    }

    err = w.Close()
    if err != nil {
        t.Error("Cannot close f")
    }
    bufferSize = oldBufferSize
}

func TestRotatingFileWriter(t *testing.T) {
    dir := filepath.Join(os.TempDir(), "test")
    path := filepath.Join(dir, "test.log")
    err := os.RemoveAll(dir)
    if err != nil {
        t.Error(err)
    }
    err = os.Mkdir(dir, 0755)
    if err != nil {
        t.Error(err)
    }

    w, err := NewRotatingFileWriter(path, 128, 2)
    if err != nil {
        t.Error(err)
    }
    stat, err := os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    bs := []byte("0123456789")
    for i := 0; i < 20; i++ {
        _, err := w.Write(bs)
        if err != nil {
            t.Error("Cannot write into w")
        }

    }

    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    stat, err = os.Stat(path + ".1")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 120 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    _, err = os.Stat(path + ".2")
    if !os.IsNotExist(err) {
        t.Error(err)
    }

    time.Sleep(flushDuration * 2)
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 80 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    // second write
    for i := 0; i < 20; i++ {
        _, err := w.Write(bs)
        if err != nil {
            t.Error("Cannot write into w")
        }
    }

    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    stat, err = os.Stat(path + ".1")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 120 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    stat, err = os.Stat(path + ".2")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 120 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    time.Sleep(flushDuration * 2)
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 40 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    err = w.Close()
    if err != nil {
        t.Error("Cannot close w")
    }
}

func TestTimedRotatingFileWriterByDate(t *testing.T) {
    dir := filepath.Join(os.TempDir(), "test")
    pathPrefix := filepath.Join(dir, "test")
    err := os.RemoveAll(dir)
    if err != nil {
        t.Error(err)
    }
    err = os.Mkdir(dir, 0755)
    if err != nil {
        t.Error(err)
    }

    tm := time.Date(2018, 11, 19, 16, 12, 34, 56, time.Local)
    var locker sync.RWMutex
    setNowFunc(func() time.Time {
        locker.RLock()
        now := tm
        locker.RUnlock()
        return now
    })
    var setNow = func(now time.Time) {
        locker.Lock()
        tm = now
        locker.Unlock()
    }

    oldNextRotateDuration := nextRotateDuration
    nextRotateDuration = func(rotateDuration RotateDuration) time.Duration {
        return flushDuration * 3
    }

    w, err := NewTimedRotatingFileWriter(pathPrefix, RotateByDate, 2)
    if err != nil {
        t.Error(err)
    }
    path := pathPrefix + "-20181119.log"
    stat, err := os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    _, err = w.Write([]byte("123"))
    if err != nil {
        t.Error("Cannot write into w")
    }
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    setNow(time.Date(2018, 11, 20, 16, 12, 34, 56, time.Local))
    time.Sleep(flushDuration * 2)
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 3 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    time.Sleep(flushDuration * 2)
    path = pathPrefix + "-20181120.log"
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    _, err = w.Write([]byte("4567"))
    if err != nil {
        t.Error("Cannot write into w")
    }
    setNow(time.Date(2018, 11, 21, 16, 12, 34, 56, time.Local))
    time.Sleep(flushDuration * 4)
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 4 {
        t.Errorf("file size are %d bytes", stat.Size())
    }
    stat, err = os.Stat(pathPrefix + "-20181121.log")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    setNow(time.Date(2018, 11, 22, 16, 12, 34, 56, time.Local))
    time.Sleep(flushDuration * 4)
    stat, err = os.Stat(pathPrefix + "-20181122.log")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }
    _, err = os.Stat(pathPrefix + "-20181119.log")
    if !os.IsNotExist(err) {
        t.Error(err)
    }

    err = w.Close()
    if err != nil {
        t.Error("Cannot close w")
    }
    setNowFunc(time.Now)
    nextRotateDuration = oldNextRotateDuration
}

func TestTimedRotatingFileWriterByHour(t *testing.T) {
    dir := filepath.Join(os.TempDir(), "test")
    pathPrefix := filepath.Join(dir, "test")
    err := os.RemoveAll(dir)
    if err != nil {
        t.Error(err)
    }
    err = os.Mkdir(dir, 0755)
    if err != nil {
        t.Error(err)
    }

    tm := time.Date(2018, 11, 19, 16, 12, 34, 56, time.Local)
    var locker sync.RWMutex
    setNowFunc(func() time.Time {
        locker.RLock()
        now := tm
        locker.RUnlock()
        return now
    })
    var setNow = func(now time.Time) {
        locker.Lock()
        tm = now
        locker.Unlock()
    }

    oldNextRotateDuration := nextRotateDuration
    nextRotateDuration = func(rotateDuration RotateDuration) time.Duration {
        return flushDuration * 3
    }

    w, err := NewTimedRotatingFileWriter(pathPrefix, RotateByHour, 2)
    if err != nil {
        t.Error(err)
    }
    path := pathPrefix + "-2018111916.log"
    stat, err := os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    _, err = w.Write([]byte("123"))
    if err != nil {
        t.Error("Cannot write into w")
    }
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    setNow(time.Date(2018, 11, 19, 17, 12, 34, 56, time.Local))
    time.Sleep(flushDuration * 2)
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 3 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    time.Sleep(flushDuration * 2)
    path = pathPrefix + "-2018111917.log"
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    _, err = w.Write([]byte("4567"))
    if err != nil {
        t.Error("Cannot write into w")
    }
    setNow(time.Date(2018, 11, 19, 18, 12, 34, 56, time.Local))
    time.Sleep(flushDuration * 4)
    stat, err = os.Stat(path)
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 4 {
        t.Errorf("file size are %d bytes", stat.Size())
    }
    stat, err = os.Stat(pathPrefix + "-2018111918.log")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }

    setNow(time.Date(2018, 11, 22, 16, 12, 34, 56, time.Local))
    time.Sleep(flushDuration * 4)
    stat, err = os.Stat(pathPrefix + "-2018112216.log")
    if err != nil {
        t.Error(err)
    }
    if stat.Size() != 0 {
        t.Errorf("file size are %d bytes", stat.Size())
    }
    _, err = os.Stat(pathPrefix + "-2018111916.log")
    if !os.IsNotExist(err) {
        t.Error(err)
    }

    err = w.Close()
    if err != nil {
        t.Error("Cannot close w")
    }
    setNowFunc(time.Now)
    nextRotateDuration = oldNextRotateDuration
}

type badWriter struct{}

func (w *badWriter) Write(p []byte) (n int, err error) {
    return 0, io.ErrShortWrite
}

func (w *badWriter) Close() error {
    return nil
}

func TestBadWriter(t *testing.T) {
    oldLogger := internalLogger
    SetInternalLogger(nil)

    l := NewLoggerWithWriter(&badWriter{})
    l.Log(InfoLevel, "", 0, "test")
    logError(errors.New("test"))

    SetInternalLogger(NewLoggerWithWriter(&badWriter{}))
    l.Log(InfoLevel, "", 0, "test")
    logError(errors.New("test"))

    SetInternalLogger(oldLogger)
    l.Close()
}

func TestNewStdoutWriter(t *testing.T) {
    l := NewStdoutWriter()
    content := "Test"
    n, err := l.Write([]byte(content))
    if err != nil {
        t.Fatal("Cannot write to Std Out")
    }
    assert.Equalf(t, n, utf8.RuneCountInString(content), "the length of content should be %v", n)
}

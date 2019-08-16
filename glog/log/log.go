package log

import (
    "runtime"

    "github.com/Tomilla/imagespider/glog"
)

var defaultLogger *glog.Logger

var (
    logFuncs = map[glog.Level](func(args ...interface{})){
        glog.DebugLevel: _debug,
        glog.InfoLevel:  _info,
        glog.WarnLevel:  _warn,
        glog.ErrorLevel: _error,
        glog.CritLevel:  _crit,
    }
    logfFuncs = map[glog.Level](func(msg string, args ...interface{})){
        glog.DebugLevel: _debugf,
        glog.InfoLevel:  _infof,
        glog.WarnLevel:  _warnf,
        glog.ErrorLevel: _errorf,
        glog.CritLevel:  _critf,
    }
    logPtrs = map[glog.Level](*func(args ...interface{})){
        glog.DebugLevel: &Debug,
        glog.InfoLevel:  &Info,
        glog.WarnLevel:  &Warn,
        glog.ErrorLevel: &Error,
        glog.CritLevel:  &Crit,
    }
    logfPtrs = map[glog.Level](*func(msg string, args ...interface{})){
        glog.DebugLevel: &Debugf,
        glog.InfoLevel:  &Infof,
        glog.WarnLevel:  &Warnf,
        glog.ErrorLevel: &Errorf,
        glog.CritLevel:  &Critf,
    }
)

// SetDefaultLogger set the logger as the defaultLogger.
// The logging functions in this package use it as their logger.
// This function should be called before using the others.
func SetDefaultLogger(l *glog.Logger) {
    defaultLogger = l

    minLevel := l.GetMinLevel()
    for level, f := range logFuncs {
        if minLevel <= level {
            *logPtrs[level] = f
        } else {
            *logPtrs[level] = nop
        }
    }
    for level, f := range logfFuncs {
        if minLevel <= level {
            *logfPtrs[level] = f
        } else {
            *logfPtrs[level] = nopf
        }
    }
}

func nop(args ...interface{})              {}
func nopf(msg string, args ...interface{}) {}

// Debug logs a _debug level message. It uses fmt.Fprint() to format args.
var Debug func(args ...interface{})

// Debugf logs a _debug level message. It uses fmt.Fprintf() to format msg and args.
var Debugf func(msg string, args ...interface{})

// Info logs a _info level message. It uses fmt.Fprint() to format args.
var Info func(args ...interface{})

// Infof logs a _info level message. It uses fmt.Fprintf() to format msg and args.
var Infof func(msg string, args ...interface{})

// Warn logs a _warning level message. It uses fmt.Fprint() to format args.
var Warn func(args ...interface{})

// Warnf logs a _warning level message. It uses fmt.Fprintf() to format msg and args.
var Warnf func(msg string, args ...interface{})

// Error logs an _error level message. It uses fmt.Fprint() to format args.
var Error func(args ...interface{})

// Errorf logs a _error level message. It uses fmt.Fprintf() to format msg and args.
var Errorf func(msg string, args ...interface{})

// Crit logs a _critical level message. It uses fmt.Fprint() to format args.
var Crit func(args ...interface{})

// Critf logs a _critical level message. It uses fmt.Fprintf() to format msg and args.
var Critf func(msg string, args ...interface{})

func _debug(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1) // deeper caller will be more expensive
    defaultLogger.Log(glog.DebugLevel, file, line, "", args...)
}

func _debugf(msg string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.DebugLevel, file, line, msg, args...)
}

func _info(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.InfoLevel, file, line, "", args...)
}

func _infof(msg string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.InfoLevel, file, line, msg, args...)
}

func _warn(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.WarnLevel, file, line, "", args...)
}

func _warnf(msg string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.WarnLevel, file, line, msg, args...)
}

func _error(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.ErrorLevel, file, line, "", args...)
}

func _errorf(msg string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.ErrorLevel, file, line, msg, args...)
}

func _crit(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.CritLevel, file, line, "", args...)
}

func _critf(msg string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    defaultLogger.Log(glog.CritLevel, file, line, msg, args...)
}

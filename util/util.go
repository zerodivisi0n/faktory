package util

import (
	"bufio"
	"bytes"
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"os"
	"reflect"
	"runtime"
	"time"
)

var (
	LogInfo  = false
	LogDebug = false
	logg     = NewLogger("info", false)
)

func Darwin() bool {
	b, _ := FileExists("/Applications")
	return b
}

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(data []byte) ([]string, error) {
	var lines []string
	scan := bufio.NewScanner(bytes.NewReader(data))
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	return lines, scan.Err()
}

//
// Logging functions
//

func InitLogger(level string) {
	logg = NewLogger(level, true)
	if level == "info" {
		LogInfo = true
	}

	if level == "debug" {
		LogInfo = true
		LogDebug = true
	}
}

func Log() Logger {
	return logg
}

func Error(msg string, err error, stack []byte) {
	logg.Error(msg, reflect.TypeOf(err).Name(), err.Error())
	if stack != nil {
		logg.Error(string(stack))
	}
}

// Uh oh, not good but not worthy of process death
func Warn(args ...interface{}) {
	logg.Warn(args...)
}

func Warnf(msg string, args ...interface{}) {
	logg.Warnf(msg, args...)
}

// Typical logging output, the default level
func Info(args ...interface{}) {
	if LogInfo {
		logg.Info(args...)
	}
}

// Typical logging output, the default level
func Infof(msg string, args ...interface{}) {
	if LogInfo {
		logg.Infof(msg, args...)
	}
}

// -l debug: Verbosity level which helps track down production issues
func Debug(args ...interface{}) {
	if LogDebug {
		logg.Debug(args...)
	}
}

// -l debug: Verbosity level which helps track down production issues
func Debugf(msg string, args ...interface{}) {
	if LogDebug {
		logg.Debugf(msg, args...)
	}
}

func RandomJid() string {
	bytes := make([]byte, 12)
	_, err := cryptorand.Read(bytes)
	if err != nil {
		mathrand.Read(bytes)
	}

	return base64.RawURLEncoding.EncodeToString(bytes)
}

const (
	// This is the canonical timestamp format used by Faktory.
	// Always UTC, lexigraphically sortable.  This is the best
	// timestamp format, accept no others.
	TimestampFormat = time.RFC3339Nano
)

func Thens(tim time.Time) string {
	return tim.UTC().Format(TimestampFormat)
}

func Nows() string {
	return time.Now().UTC().Format(TimestampFormat)
}

func ParseTime(str string) (time.Time, error) {
	return time.Parse(TimestampFormat, str)
}

/*
 * Gather a backtrace for the caller.
 * Return a slice of up to N stack frames.
 */
func Backtrace(size int) []string {
	pc := make([]uintptr, size)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return []string{}
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	str := make([]string, size)
	count := 0

	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for i := 0; i < size; i++ {
		frame, more := frames.Next()
		str[i] = fmt.Sprintf("in %s:%d %s", frame.File, frame.Line, frame.Function)
		count += 1
		if !more {
			break
		}
	}

	return str[0:count]
}

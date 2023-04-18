package log

import (
	"bytes"
	"fmt"
	"github.com/liuyp5181/base/signal"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	fileName = "service"
)

const (
	consoleMode = 0
	fileMode    = 1
)

const (
	debugLv = iota
	infoLv
	warningLv
	errorLv
	fatalLv
)

var logName = map[int]string{
	debugLv:   "DEBUG",
	infoLv:    "INFO",
	warningLv: "WARNING",
	errorLv:   "ERROR",
	fatalLv:   "FATAL",
}

type Config struct {
	Path    string `mapstructure:"path"`
	Name    string `mapstructure:"name"`
	Level   int    `mapstructure:"level"`
	MaxSize int64  `mapstructure:"max_size"`
	MaxAge  int    `mapstructure:"max_age"`
}

type Logger struct {
	cfg      *Config
	f        *os.File
	ch       chan string
	closeCh  chan struct{}
	wg       sync.WaitGroup
	currTime string
	mode     int
}

var logger Logger

func init() {
	logger.ch = make(chan string, 100)
	logger.closeCh = make(chan struct{})
	logger.wg.Add(1)
	logger.cfg = &Config{
		Level:   debugLv,
		MaxSize: 0,
		MaxAge:  0,
	}
	logger.f = os.Stderr
	logger.mode = consoleMode
	go watch()

	signal.RegisterClose(func() {
		close(logger.closeCh)
		logger.wg.Wait()
	})
}

func Init(cfg *Config) error {
	cfg.MaxSize *= 1024 * 1024
	logger.cfg = cfg
	if len(cfg.Name) == 0 {
		cfg.Name = fileName
	}
	tn := time.Now().Format("2006-01-02")
	fn := fmt.Sprintf("%s/%s-%s.%d.log", cfg.Path, cfg.Name, tn, 1)
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	logger.currTime = tn
	logger.f = f
	if isNewFile() {
		logger.f.Close()
		f, _ := openFile(logger.cfg)
		logger.f = f
	}
	logger.mode = fileMode

	return nil
}

func isNewFile() bool {
	if logger.f == nil {
		return true
	}

	if logger.currTime != time.Now().Format("2006-01-02") {
		return true
	}

	fi, _ := logger.f.Stat()
	println(fi.Size(), logger.cfg.MaxSize)
	if logger.cfg.MaxSize > 0 && fi.Size() >= logger.cfg.MaxSize {
		return true
	}

	return false
}

func openFile(cfg *Config) (*os.File, error) {
	var fn string
	var i int
	var tn string
	for {
		i++
		tn = time.Now().Format("2006-01-02")
		fn = fmt.Sprintf("%s/%s-%s.%d.log", cfg.Path, cfg.Name, tn, i)
		_, err := os.Stat(fn)
		if err != nil && os.IsNotExist(err) {
			// "File not exist"
			break
		}
		// "File exist"
	}

	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	logger.currTime = tn
	return f, nil
}

func write(s string) {
	defer logger.f.WriteString(s)

	if logger.mode == consoleMode {
		return
	}

	if isNewFile() {
		logger.f.Close()
		f, _ := openFile(logger.cfg)
		logger.f = f
	}
}

func watch() {
	var logs = make([]string, 0, 10)
	var s string
	var ticker = time.NewTicker(800 * time.Millisecond)
	defer logger.wg.Done()
	for {
		select {
		case s = <-logger.ch:
			logs = append(logs, s)
			if len(logs) < 10 {
				continue
			}
			write(strings.Join(logs, ""))
			logs = make([]string, 0, 10)
		case <-ticker.C:
			if len(logs) == 0 {
				continue
			}
			write(strings.Join(logs, ""))
			logs = make([]string, 0, 10)
		case <-logger.closeCh:
			if len(logs) == 0 {
				return
			}
			write(strings.Join(logs, ""))
			logs = make([]string, 0, 10)
			return
		}
	}
}

func header(lv int) string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.SplitAfter(file, "/")
		l := len(slash)
		if l >= 2 {
			file = slash[l-2] + slash[l-1]
		}
	}
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	b := buf[10:n]
	l := bytes.IndexByte(b, ' ')
	return fmt.Sprintf("%s %s:%d:g%s [%s] ", time.Now().Format("2006-01-02 15:04:05"), file, line, b[:l], logName[lv])
}

func output(lv int, args ...interface{}) {
	if lv < logger.cfg.Level {
		return
	}
	logger.ch <- header(lv) + fmt.Sprintln(args...)
}

func outputf(lv int, format string, args ...interface{}) {
	if lv < logger.cfg.Level {
		return
	}
	logger.ch <- header(lv) + fmt.Sprintf(format, args...) + "\n"
}

func Debug(args ...interface{}) {
	output(debugLv, args...)
}
func Debugf(format string, args ...interface{}) {
	outputf(debugLv, format, args...)
}

func Info(args ...interface{}) {
	output(infoLv, args...)
}
func Infof(format string, args ...interface{}) {
	outputf(infoLv, format, args...)
}

func Warning(args ...interface{}) {
	output(warningLv, args...)
}
func Warningf(format string, args ...interface{}) {
	outputf(warningLv, format, args...)
}

func Error(args ...interface{}) {
	output(errorLv, args...)
}
func Errorf(format string, args ...interface{}) {
	outputf(errorLv, format, args...)
}

func Fatal(args ...interface{}) {
	output(fatalLv, args...)
	close(logger.closeCh)
	logger.wg.Wait()
	os.Exit(1114)
}
func Fatalf(format string, args ...interface{}) {
	outputf(fatalLv, format, args...)
	close(logger.closeCh)
	logger.wg.Wait()
	os.Exit(1124)
}

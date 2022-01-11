package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/valyala/fasttemplate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
)

type (
	CLogger struct {
		prefix     string
		level      uint32
		skip       int
		output     io.Writer
		template   *fasttemplate.Template
		levels     []string
		color      *color.Color
		bufferPool sync.Pool
		mutex      sync.Mutex
	}

	Lvl uint8

	JSON map[string]interface{}
)

const (
	DEBUG Lvl = iota + 1
	INFO
	WARN
	ERROR
	OFF
	panicLevel
	fatalLevel
)

func ParseLvl(lvl string) Lvl {
	levels := map[string]Lvl{
		"DEBUG": DEBUG,
		"INFO":  INFO,
		"WARN":  WARN,
		"ERROR": ERROR,
		"FATAL": fatalLevel,
	}
	return levels[lvl]
}

var (
	global        = New("-")
	defaultHeader = `{"time":"${time_rfc3339_nano}","level":"${level}", "traceId":"${traceId}", "spanId":"${spanId}", "prefix":"${prefix}",` +
		`"file":"${short_file}","line":"${line}"}`
)

func init() {
	global.skip = 3
}

func New(prefix string) (l *CLogger) {
	l = &CLogger{
		level:    uint32(INFO),
		skip:     2,
		prefix:   prefix,
		template: l.newTemplate(defaultHeader),
		color:    color.New(),
		bufferPool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 256))
			},
		},
	}
	l.initLevels()
	l.SetOutput(output())

	return
}

func output() io.Writer {
	return colorable.NewColorableStdout()
}

func (l *CLogger) initLevels() {
	l.levels = []string{
		"-",
		l.color.Blue("DEBUG"),
		l.color.Green("INFO"),
		l.color.Yellow("WARN"),
		l.color.Red("ERROR"),
		"",
		l.color.Yellow("PANIC", color.U),
		l.color.Red("FATAL", color.U),
	}
}

func (l *CLogger) newTemplate(format string) *fasttemplate.Template {
	return fasttemplate.New(format, "${", "}")
}

func (l *CLogger) DisableColor() {
	l.color.Disable()
	l.initLevels()
}

func (l *CLogger) EnableColor() {
	l.color.Enable()
	l.initLevels()
}

func (l *CLogger) Prefix() string {
	return l.prefix
}

func (l *CLogger) SetPrefix(p string) {
	l.prefix = p
}

func (l *CLogger) Level() Lvl {
	return Lvl(atomic.LoadUint32(&l.level))
}

func (l *CLogger) SetLevel(level Lvl) {
	atomic.StoreUint32(&l.level, uint32(level))
}

func (l *CLogger) Output() io.Writer {
	return l.output
}

func (l *CLogger) SetOutput(w io.Writer) {
	l.output = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		l.DisableColor()
	}
}

func (l *CLogger) Color() *color.Color {
	return l.color
}

func (l *CLogger) SetHeader(h string) {
	l.template = l.newTemplate(h)
}

func (l *CLogger) Print(i ...interface{}) {
	l.log(0, "", i...)
}

func (l *CLogger) CPrint(ctx echo.Context, i ...interface{}) {
	l.logWithContext(0, "", ctx, i...)
}

func (l *CLogger) Printf(format string, args ...interface{}) {
	l.log(0, format, args...)
}

func (l *CLogger) CPrintf(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(0, format, ctx, args...)
}

func (l *CLogger) Printj(j JSON) {
	l.log(0, "json", j)
}

func (l *CLogger) CPrintj(ctx echo.Context, j JSON) {
	l.logWithContext(0, "json", ctx, j)
}

func (l *CLogger) Debug(i ...interface{}) {
	l.log(DEBUG, "", i...)
}

func (l *CLogger) CDebug(ctx echo.Context, i ...interface{}) {
	l.logWithContext(DEBUG, "", ctx, i...)
}

func (l *CLogger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *CLogger) CDebugf(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(DEBUG, format, ctx, args...)
}

func (l *CLogger) Debugj(j JSON) {
	l.log(DEBUG, "json", j)
}

func (l *CLogger) CDebugj(ctx echo.Context, j JSON) {
	l.logWithContext(DEBUG, "json", ctx, j)
}

func (l *CLogger) Info(i ...interface{}) {
	l.log(INFO, "", i...)
}

func (l *CLogger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *CLogger) Infoj(j JSON) {
	l.log(INFO, "json", j)
}

func (l *CLogger) CInfo(ctx echo.Context, i ...interface{}) {
	l.logWithContext(INFO, "", ctx, i...)
}

func (l *CLogger) CInfof(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(INFO, format, ctx, args...)
}

func (l *CLogger) CInfoj(ctx echo.Context, j JSON) {
	l.logWithContext(INFO, "json", ctx, j)
}

func (l *CLogger) Warn(i ...interface{}) {
	l.log(WARN, "", i...)
}

func (l *CLogger) Warnf(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *CLogger) Warnj(j JSON) {
	l.log(WARN, "json", j)
}

func (l *CLogger) CWarn(ctx echo.Context, i ...interface{}) {
	l.logWithContext(WARN, "", ctx, i...)
}

func (l *CLogger) CWarnf(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(WARN, format, ctx, args...)
}

func (l *CLogger) CWarnj(ctx echo.Context, j JSON) {
	l.logWithContext(WARN, "json", ctx, j)
}

func (l *CLogger) Error(i ...interface{}) {
	l.log(ERROR, "", i...)
}

func (l *CLogger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *CLogger) Errorj(j JSON) {
	l.log(ERROR, "json", j)
}

func (l *CLogger) CError(ctx echo.Context, i ...interface{}) {
	l.logWithContext(ERROR, "", ctx, i...)
}

func (l *CLogger) CErrorf(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(ERROR, format, ctx, args...)
}

func (l *CLogger) CErrorj(ctx echo.Context, j JSON) {
	l.logWithContext(ERROR, "json", ctx, j)
}

func (l *CLogger) Fatal(i ...interface{}) {
	l.log(fatalLevel, "", i...)
	os.Exit(1)
}

func (l *CLogger) Fatalf(format string, args ...interface{}) {
	l.log(fatalLevel, format, args...)
	os.Exit(1)
}

func (l *CLogger) Fatalj(j JSON) {
	l.log(fatalLevel, "json", j)
	os.Exit(1)
}

func (l *CLogger) CFatal(ctx echo.Context, i ...interface{}) {
	l.logWithContext(fatalLevel, "", ctx, i...)
	os.Exit(1)
}

func (l *CLogger) CFatalf(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(fatalLevel, format, ctx, args...)
	os.Exit(1)
}

func (l *CLogger) CFatalj(ctx echo.Context, j JSON) {
	l.logWithContext(fatalLevel, "json", ctx, j)
	os.Exit(1)
}

func (l *CLogger) Panic(i ...interface{}) {
	l.log(panicLevel, "", i...)
	panic(fmt.Sprint(i...))
}

func (l *CLogger) Panicf(format string, args ...interface{}) {
	l.log(panicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *CLogger) Panicj(j JSON) {
	l.log(panicLevel, "json", j)
	panic(j)
}

func (l *CLogger) CPanic(ctx echo.Context, i ...interface{}) {
	l.logWithContext(panicLevel, "", ctx, i...)
	panic(fmt.Sprint(i...))
}

func (l *CLogger) CPanicf(ctx echo.Context, format string, args ...interface{}) {
	l.logWithContext(panicLevel, format, ctx, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *CLogger) CPanicj(ctx echo.Context, j JSON) {
	l.logWithContext(panicLevel, "json", ctx, j)
	panic(j)
}

func DisableColor() {
	global.DisableColor()
}

func EnableColor() {
	global.EnableColor()
}

func Prefix() string {
	return global.Prefix()
}

func SetPrefix(p string) {
	global.SetPrefix(p)
}

func Level() Lvl {
	return global.Level()
}

func SetLevel(level Lvl) {
	global.SetLevel(level)
}

func Output() io.Writer {
	return global.Output()
}

func SetOutput(w io.Writer) {
	global.SetOutput(w)
}

func SetHeader(h string) {
	global.SetHeader(h)
}

func Print(i ...interface{}) {
	global.Print(i...)
}

func Printf(format string, args ...interface{}) {
	global.Printf(format, args...)
}

func Printj(j JSON) {
	global.Printj(j)
}

func CPrint(ctx echo.Context, i ...interface{}) {
	global.CPrint(ctx, i...)
}

func CPrintf(ctx echo.Context, format string, args ...interface{}) {
	global.CPrintf(ctx, format, args...)
}

func CPrintj(ctx echo.Context, j JSON) {
	global.CPrintj(ctx, j)
}

func Debug(i ...interface{}) {
	global.Debug(i...)
}

func Debugf(format string, args ...interface{}) {
	global.Debugf(format, args...)
}

func Debugj(j JSON) {
	global.Debugj(j)
}

func CDebug(ctx echo.Context, i ...interface{}) {
	global.CDebug(ctx, i...)
}

func CDebugf(ctx echo.Context, format string, args ...interface{}) {
	global.CDebugf(ctx, format, args...)
}

func CDebugj(ctx echo.Context, j JSON) {
	global.CDebugj(ctx, j)
}

func Info(i ...interface{}) {
	global.Info(i...)
}

func Infof(format string, args ...interface{}) {
	global.Infof(format, args...)
}

func Infoj(j JSON) {
	global.Infoj(j)
}

func CInfo(ctx echo.Context, i ...interface{}) {
	global.CInfo(ctx, i...)
}

func CInfof(ctx echo.Context, format string, args ...interface{}) {
	global.CInfof(ctx, format, args...)
}

func CInfoj(ctx echo.Context, j JSON) {
	global.CInfoj(ctx, j)
}

func Warn(i ...interface{}) {
	global.Warn(i...)
}

func Warnf(format string, args ...interface{}) {
	global.Warnf(format, args...)
}

func Warnj(j JSON) {
	global.Warnj(j)
}

func CWarn(ctx echo.Context, i ...interface{}) {
	global.CWarn(ctx, i...)
}

func CWarnf(ctx echo.Context, format string, args ...interface{}) {
	global.CWarnf(ctx, format, args...)
}

func CWarnj(ctx echo.Context, j JSON) {
	global.CWarnj(ctx, j)
}

func Error(i ...interface{}) {
	global.Error(i...)
}

func Errorf(format string, args ...interface{}) {
	global.Errorf(format, args...)
}

func Errorj(j JSON) {
	global.Errorj(j)
}

func CError(ctx echo.Context, i ...interface{}) {
	global.CError(ctx, i...)
}

func CErrorf(ctx echo.Context, format string, args ...interface{}) {
	global.CErrorf(ctx, format, args...)
}

func CErrorj(ctx echo.Context, j JSON) {
	global.CErrorj(ctx, j)
}

func Fatal(i ...interface{}) {
	global.Fatal(i...)
}

func Fatalf(format string, args ...interface{}) {
	global.Fatalf(format, args...)
}

func Fatalj(j JSON) {
	global.Fatalj(j)
}

func CFatal(ctx echo.Context, i ...interface{}) {
	global.CFatal(ctx, i...)
}

func CFatalf(ctx echo.Context, format string, args ...interface{}) {
	global.CFatalf(ctx, format, args...)
}

func CFatalj(ctx echo.Context, j JSON) {
	global.CFatalj(ctx, j)
}

func Panic(i ...interface{}) {
	global.Panic(i...)
}

func Panicf(format string, args ...interface{}) {
	global.Panicf(format, args...)
}

func Panicj(j JSON) {
	global.Panicj(j)
}

func CPanic(ctx echo.Context, i ...interface{}) {
	global.CPanic(ctx, i...)
}

func CPanicf(ctx echo.Context, format string, args ...interface{}) {
	global.CPanicf(ctx, format, args...)
}

func CPanicj(ctx echo.Context, j JSON) {
	global.CPanicj(ctx, j)
}

func (l *CLogger) logWithContext(level Lvl, format string, ctx echo.Context, args ...interface{}) {
	if level >= l.Level() || level == 0 {
		buf := l.bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer l.bufferPool.Put(buf)

		_, file, line, _ := runtime.Caller(l.skip)
		message := ""

		if format == "" {
			message = fmt.Sprint(args...)
		} else if format == "json" {
			b, err := json.Marshal(args[0])
			if err != nil {
				panic(err)
			}
			message = string(b)
		} else {
			message = fmt.Sprintf(format, args...)
		}

		_, err := l.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "time_rfc3339":
				return w.Write([]byte(time.Now().Format(time.RFC3339)))
			case "time_rfc3339_nano":
				return w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
			case "level":
				return w.Write([]byte(l.levels[level]))
			case "prefix":
				return w.Write([]byte(l.prefix))
			case "long_file":
				return w.Write([]byte(file))
			case "short_file":
				return w.Write([]byte(path.Base(file)))
			case "line":
				return w.Write([]byte(strconv.Itoa(line)))
			case "traceId":
				return w.Write([]byte(GetTraceId(ctx)))
			case "spanId":
				return w.Write([]byte(GetSpanId(ctx)))
			}
			return 0, nil
		})
		if err == nil {
			s := buf.String()
			i := buf.Len() - 1
			if s[i] == '}' {
				// JSON header
				buf.Truncate(i)
				buf.WriteByte(',')
				if format == "json" {
					buf.WriteString(message[1:])
				} else {
					buf.WriteString(`"message":`)
					buf.WriteString(strconv.Quote(message))
					buf.WriteString(`}`)
				}
			} else {
				// Text header
				buf.WriteByte(' ')
				buf.WriteString(message)
			}
			buf.WriteByte('\n')
			l.mutex.Lock()
			defer l.mutex.Unlock()
			_, _ = l.output.Write(buf.Bytes())
		}
	}
}

func (l *CLogger) log(level Lvl, format string, args ...interface{}) {
	if level >= l.Level() || level == 0 {
		buf := l.bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer l.bufferPool.Put(buf)
		_, file, line, _ := runtime.Caller(l.skip)
		message := ""

		if format == "" {
			message = fmt.Sprint(args...)
		} else if format == "json" {
			b, err := json.Marshal(args[0])
			if err != nil {
				panic(err)
			}
			message = string(b)
		} else {
			message = fmt.Sprintf(format, args...)
		}

		_, err := l.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "time_rfc3339":
				return w.Write([]byte(time.Now().Format(time.RFC3339)))
			case "time_rfc3339_nano":
				return w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
			case "level":
				return w.Write([]byte(l.levels[level]))
			case "prefix":
				return w.Write([]byte(l.prefix))
			case "long_file":
				return w.Write([]byte(file))
			case "short_file":
				return w.Write([]byte(path.Base(file)))
			case "line":
				return w.Write([]byte(strconv.Itoa(line)))
			}
			return 0, nil
		})

		if err == nil {
			s := buf.String()
			i := buf.Len() - 1
			if s[i] == '}' {
				// JSON header
				buf.Truncate(i)
				buf.WriteByte(',')
				if format == "json" {
					buf.WriteString(message[1:])
				} else {
					buf.WriteString(`"message":`)
					buf.WriteString(strconv.Quote(message))
					buf.WriteString(`}`)
				}
			} else {
				// Text header
				buf.WriteByte(' ')
				buf.WriteString(message)
			}
			buf.WriteByte('\n')
			l.mutex.Lock()
			defer l.mutex.Unlock()
			_, _ = l.output.Write(buf.Bytes())
		}
	}
}

package middlewares

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/govm/entity"
	"github.com/sivaosorg/govm/logger"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *gin.Context, err any)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func (m *MiddlewareManager) Recovery() gin.HandlerFunc {
	return m.RecoveryWithWriter(gin.DefaultErrorWriter)
}

// CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.
func (m *MiddlewareManager) CustomRecovery(handle RecoveryFunc) gin.HandlerFunc {
	return m.RecoveryWithWriter(gin.DefaultErrorWriter, handle)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func (m *MiddlewareManager) RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) gin.HandlerFunc {
	if len(recovery) > 0 {
		return m.CustomRecoveryWithWriter(out, recovery[0])
	}
	return m.CustomRecoveryWithWriter(out, m.defaultHandleRecovery)
}

// CustomRecoveryWithWriter returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it.
func (m *MiddlewareManager) CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				stack := stack(3)
				logger.Debugf("RecoveryWriter: Panic recovered: %v %s", err, string(stack))
				if brokenPipe {
					c.Abort()
				} else {
					handle(c, err)
				}
			}
		}()
		c.Next()
	}
}

func (m *MiddlewareManager) defaultHandleRecovery(c *gin.Context, err any) {
	response := entity.NewResponseEntity().InternalServerError(http.StatusText(http.StatusInternalServerError), nil)
	response.SetErrors(fmt.Sprint(err))
	m.notification(c, err, response.StatusCode)
	c.JSON(response.StatusCode, response)
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

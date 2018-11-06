package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Handler defines where and how log records are written.
// A Logger prints its log records by writing to a Handler.
// Handlers are composable, providing you great flexibility in combining
// them to achieve the logging structure that suits your applications.
type Handler interface {
	Log(r *Record) error
}

// FuncHandler returns a Handler that logs records with the given
// function.
func FuncHandler(fn func(r *Record) error) Handler {
	return funcHandler(fn)
}

type funcHandler func(r *Record) error

func (h funcHandler) Log(r *Record) error {
	return h(r)
}

// StreamHandler writes log records to an io.Writer
// with the given format. StreamHandler can be used
// to easily begin writing log records to other
// outputs.
//
// StreamHandler wraps itself with LazyHandler and SyncHandler
// to evaluate Lazy objects and perform safe concurrent writes.
func StreamHandler(wr io.Writer, fmtr Format) Handler {
	h := FuncHandler(func(r *Record) error {
		_, err := wr.Write(fmtr.Format(r))
		return err
	})
	//return LazyHandler(SyncHandler(h))
	return SyncHandler(h)
}

// SyncHandler can be wrapped around a handler to guarantee that
// only a single Log operation can proceed at a time. It's necessary
// for thread-safe concurrent writes.
func SyncHandler(h Handler) Handler {
	var mu sync.Mutex
	return FuncHandler(func(r *Record) error {
		defer mu.Unlock()
		mu.Lock()
		return h.Log(r)
	})
}

// FileHandler returns a handler which writes log records to the give file
// using the given format. If the path
// already exists, FileHandler will append to the given file. If it does not,
// FileHandler will create the file with mode 0644.
func FileHandler(path string, fmtr Format) (Handler, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return closingHandler{f, StreamHandler(f, fmtr)}, nil
}

type closingHandler struct {
	io.WriteCloser
	Handler
}

func (h *closingHandler) Close() error {
	return h.WriteCloser.Close()
}

// prepFile opens the log file at the given path, and cuts off the invalid part
// from the end, because the previous execution could have been finished by interruption.
// Assumes that every line ended by '\n' contains a valid log record.
func prepFile(path string) (*countingWriter, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(-1, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1)
	var cut int64
	for {
		if _, err := f.Read(buf); err != nil {
			return nil, err
		}
		if buf[0] == '\n' {
			break
		}
		if _, err = f.Seek(-1, io.SeekCurrent); err != nil {
			return nil, err
		}
		cut++
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	ns := fi.Size() - cut
	if err = f.Truncate(ns); err != nil {
		return nil, err
	}
	return &countingWriter{w: f, count: uint(ns)}, nil
}

// RotatingFileHandler returns a handler which writes log records to file chunks
// at the given path. When a file's size reaches the limit, the handler creates
// a new file named after the timestamp of the first log record it will contain.
func RotatingFileHandler(path string, limit uint, formatter Format) (Handler, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`\.log$`)
	last := len(files) - 1
	for last >= 0 && (!files[last].Mode().IsRegular() || !re.MatchString(files[last].Name())) {
		last--
	}
	var counter *countingWriter
	if last >= 0 && files[last].Size() < int64(limit) {
		// Open the last file, and continue to write into it until it's size reaches the limit.
		if counter, err = prepFile(filepath.Join(path, files[last].Name())); err != nil {
			return nil, err
		}
	}
	if counter == nil {
		counter = new(countingWriter)
	}
	h := StreamHandler(counter, formatter)

	return FuncHandler(func(r *Record) error {
		if counter.count > limit {
			counter.Close()
			counter.w = nil
		}
		if counter.w == nil {
			f, err := os.OpenFile(
				filepath.Join(path, fmt.Sprintf("%s.log", strings.Replace(r.Time.Format("060102150405.00"), ".", "", 1))),
				os.O_CREATE|os.O_APPEND|os.O_WRONLY,
				0600,
			)
			if err != nil {
				return err
			}
			counter.w = f
			counter.count = 0
		}
		return h.Log(r)
	}), nil
}

// NetHandler opens a socket to the given address and writes records
// over the connection.
func NetHandler(network, addr string, fmtr Format) (Handler, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return closingHandler{conn, StreamHandler(conn, fmtr)}, nil
}

//countingWriter wraps a WriteCloser object in order to count the written bytes.
type countingWriter struct {
	w     io.WriteCloser
	count uint
}

// Write increments the byte counter by the number of bytes written.
// Implements the WriteCloser interface.
func (w *countingWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	w.count += uint(n)
	return n, err
}

// Close implements the WriteCloser interface.
func (w *countingWriter) Close() error {
	return w.w.Close()
}

//func LazyHandler(h Handler) Handler{
//	return FuncHandler(func(r *Record)error{
//		hadErr :=false
//		for i:=1,i<len(r.ctx);i+=2{
//
//		}
//	})
//}

// DiscardHandler reports success for all writes but does nothing.
// It is useful for dynamically disabling logging at runtime via
// a Logger's SetHandler method.
func DiscardHandler() Handler {
	return FuncHandler(func(r *Record) error {
		return nil
	})
}

// LvlFilterHandler returns a Handler that only writes
// records which are less than the given verbosity
// level to the wrapped Handler. For example, to only
// log Error/Crit records:
//
//     log.LvlFilterHandler(log.LvlError, log.StdoutHandler)
//
func LvlFilterHandler(maxLvl Lvl, h Handler) Handler {
	return FilterHandler(func(r *Record) (pass bool) {
		return r.Lvl >= maxLvl
	}, h)
}

// FilterHandler returns a Handler that only writes records to the
// wrapped Handler if the given function evaluates true. For example,
// to only log records where the 'err' key is not nil:
//
//    logger.SetHandler(FilterHandler(func(r *Record) bool {
//        for i := 0; i < len(r.Ctx); i += 2 {
//            if r.Ctx[i] == "err" {
//                return r.Ctx[i+1] != nil
//            }
//        }
//        return false
//    }, h))
//
func FilterHandler(fn func(r *Record) bool, h Handler) Handler {
	return FuncHandler(func(r *Record) error {
		if fn(r) {
			return h.Log(r)
		}
		return nil
	})
}

// MultiHandler dispatches any write to each of its handlers.
// This is useful for writing different types of log information
// to different locations. For example, to log to a file and
// standard error:
//
//     log.MultiHandler(
//         log.Must.FileHandler("/var/log/app.log", log.LogfmtFormat()),
//         log.StderrHandler)
//
func MultiHandler(hs ...Handler) Handler {
	return FuncHandler(func(r *Record) error {
		for _, h := range hs {
			// what to do about failures?
			h.Log(r)
		}
		return nil
	})
}

// FailoverHandler writes all log records to the first handler
// specified, but will failover and write to the second handler if
// the first handler has failed, and so on for all handlers specified.
// For example you might want to log to a network socket, but failover
// to writing to a file if the network fails, and then to
// standard out if the file write fails:
//
//     log.FailoverHandler(
//         log.Must.NetHandler("tcp", ":9090", log.JSONFormat()),
//         log.Must.FileHandler("/var/log/app.log", log.LogfmtFormat()),
//         log.StdoutHandler)
//
// All writes that do not go to the first handler will add context with keys of
// the form "failover_err_{idx}" which explain the error encountered while
// trying to write to the handlers before them in the list.
func FailoverHandler(hs ...Handler) Handler {
	return FuncHandler(func(r *Record) error {
		var err error
		for i, h := range hs {
			err = h.Log(r)
			if err == nil {
				return nil
			}
			r.Msg = fmt.Sprintf("failover_err_%d : ", i) + err.Error() + "   " + r.Msg
		}
		return err
	})
}

// ChannelHandler writes all records to the given channel.
// It blocks if the channel is full. Useful for async processing
// of log messages, it's used by BufferedHandler.
func ChannelHandler(recs chan<- *Record) Handler {
	return FuncHandler(func(r *Record) error {
		recs <- r
		return nil
	})
}

// BufferedHandler writes all records to a buffered
// channel of the given size which flushes into the wrapped
// handler whenever it is available for writing. Since these
// writes happen asynchronously, all writes to a BufferedHandler
// never return an error and any errors from the wrapped handler are ignored.
func BufferedHandler(bufSize int, h Handler) Handler {
	recs := make(chan *Record, bufSize)
	go func() {
		for m := range recs {
			_ = h.Log(m)
		}
	}()
	return ChannelHandler(recs)
}

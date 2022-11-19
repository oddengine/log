package log

import (
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/oddcancer/log/utils"
)

// Schedule Types.
const (
	SCHEDULE_DAILY    = "daily"
	SCHEDULE_DURATION = "duration"
)

// DefaultWriterConstraints dictionary is used to describe a set of rotatable writer.
type DefaultWriterConstraints struct {
	Directory string
	FileName  string
	Level     string
	Rotation  struct {
		MaxSize  int64
		Schedule struct {
			Type     string
			Duration string
		}
		History int
	}
}

// DefaultWriter represents a rotatable writer.
type DefaultWriter struct {
	constraints *DefaultWriterConstraints
	mtx         sync.Mutex
	files       []string
	fd          io.WriteCloser
	ticker      *time.Ticker
	size        int64
}

// Init this class.
func (me *DefaultWriter) Init(constraints *DefaultWriterConstraints) *DefaultWriter {
	me.constraints = constraints

	err := utils.MkdirAll(constraints.Directory)
	if err != nil {
		panic(err)
	}
	err = me.readdir()
	if err != nil {
		panic(err)
	}
	err = me.rotate()
	if err != nil {
		panic(err)
	}
	return me
}

func (me *DefaultWriter) readdir() error {
	arr, err := os.ReadDir(me.constraints.Directory)
	if err != nil {
		Errorf("Failed to ReadDir: %v", err)
		return err
	}

	for _, entry := range arr {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			Warnf("Failed to get file info: %v", err)
			continue
		}
		me.files = append(me.files, info.Name())
	}
	sort.Strings(me.files)
	return nil
}

// Write writes len(p) bytes from p to the underlying data stream.
func (me *DefaultWriter) Write(p []byte) (int, error) {
	me.mtx.Lock()
	defer me.mtx.Unlock()

	if me.constraints.Rotation.MaxSize > 0 && me.size+int64(len(p)) >= me.constraints.Rotation.MaxSize {
		err := me.rotate()
		if err != nil {
			Errorf("Failed to rotate log: %s", err)
			return 0, err
		}
	}

	n, err := me.fd.Write(p)
	if err != nil {
		return n, err
	}
	me.size += int64(n)
	return n, nil
}

func (me *DefaultWriter) rotate() error {
	var (
		delay time.Duration
	)

	// Close current file.
	if me.fd != nil {
		me.fd.Close()
	}

	// Remove history files.
	if me.constraints.Rotation.History > 0 {
		for len(me.files) >= me.constraints.Rotation.History {
			name := me.files[0]
			me.files = me.files[1:]

			err := os.Remove(me.constraints.Directory + name)
			if err != nil {
				Errorf("Failed to remove log: %s", err)
			}
		}
	}

	// Create a new file.
	now := time.Now()
	name := now.Format(me.constraints.FileName)

	f, err := utils.OpenFile(me.constraints.Directory+name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		Errorf("Failed to create log: %s", err)
		return err
	}
	me.files = append(me.files, name)
	me.fd = f
	me.size = 0
	Debugf(0, "New log: file=%s", me.constraints.Directory+name)

	// Start ticker.
	if me.ticker != nil {
		me.ticker.Stop()
		me.ticker = nil
	}

	switch me.constraints.Rotation.Schedule.Type {
	case SCHEDULE_DAILY:
		t, err := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02 ")+me.constraints.Rotation.Schedule.Duration, time.Local)
		if err != nil {
			Errorf("Failed to parse time: %s", err)
			return err
		}
		if now.After(t) {
			t = t.Add(24 * time.Hour)
		}
		delay = time.Until(t)
	case SCHEDULE_DURATION:
		d, err := time.ParseDuration(me.constraints.Rotation.Schedule.Duration)
		if err != nil {
			Errorf("Failed to parse duration: %s", err)
			return err
		}
		delay = d
	}

	if delay > 0 {
		Debugf(0, "About to rotate logger: delay=%dns", delay)
		me.ticker = time.NewTicker(delay)
		go me.wait()
	}
	return nil
}

func (me *DefaultWriter) wait() {
	<-me.ticker.C

	me.mtx.Lock()
	defer me.mtx.Unlock()

	err := me.rotate()
	if err != nil {
		Errorf("Failed to rotate log: %s", err)
		return
	}
}

package workdir

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (wd WorkDir) Watch(ctx context.Context, w *Sensor, fn ModifiedFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.scanForChanges(string(wd))
			if len(w.modified) > 0 {
				fn(wd, w.modified...)
				// Reset modified files, should not leak memory as
				// it's only strings
				w.modified = w.modified[:0:0]
				w.Last = time.Now()
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

type ModifiedFunc func(wd WorkDir, modified ...string)

func NewSensor() *Sensor {
	return &Sensor{
		Last:     time.Now(),
		modified: make([]string, 0),
		ignore:   []string{"#", ".git/", "vendor/"},
	}
}

type Sensor struct {
	Recursive bool
	Last      time.Time
	modified  []string
	ignore    []string
}

func (s *Sensor) scanForChanges(root string) {
	filepath.Walk(root, s.visit)
}

// Ignore returns true if the file should be ignored
func (w *Sensor) Ignore(path string, f os.FileInfo) bool {
	if f == nil { // if directory has been removed
		return true
	}
	if f.IsDir() {
		return true
	}
	for _, thing := range w.ignore {
		if strings.Contains(path, thing) {
			return true
		}
	}
	return false
}

func (w *Sensor) visit(path string, f os.FileInfo, err error) error {
	if w.Ignore(path, f) {
		return nil
	}
	if f.ModTime().After(w.Last) {
		w.modified = append(w.modified, path)
	}
	if w.Recursive {
		return nil
	}
	return filepath.SkipDir
}

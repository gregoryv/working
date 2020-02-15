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

func NewSensor(recursive bool) *Sensor {
	return &Sensor{
		Last:      time.Now(),
		modified:  make([]string, 0),
		recursive: recursive,
	}
}

type Sensor struct {
	Last      time.Time
	modified  []string
	recursive bool
}

func (s *Sensor) scanForChanges(root string) {
	filepath.Walk(root, s.visit)
}

var alwaysIgnore = []string{"#", ".git/", "vendor/"}

func (w *Sensor) ignore(path string, f os.FileInfo) bool {
	if f == nil { // if directory has been removed
		return true
	}
	if f.IsDir() {
		return true
	}
	for _, thing := range alwaysIgnore {
		if strings.Contains(path, thing) {
			return true
		}
	}
	return false
}

func (w *Sensor) visit(path string, f os.FileInfo, err error) error {
	if w.ignore(path, f) {
		return nil
	}
	if f.ModTime().After(w.Last) {
		w.modified = append(w.modified, path)
	}
	if w.recursive {
		return nil
	}
	return filepath.SkipDir
}

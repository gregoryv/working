package working

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (d *Directory) Watch(ctx context.Context, w *Sensor) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(w.Pause):
			w.scanForChanges(d.Path())
			if len(w.modified) > 0 {
				w.React(d, w.modified...)
				// Reset modified files, should not leak memory as
				// it's only strings
				w.modified = w.modified[:0:0]
				w.Last = time.Now()
			}
		}
	}
}

type ModifiedFunc func(d *Directory, modified ...string)

// NewSensor returns a sensor with 1s delay and no reaction func.
// Set React.
func NewSensor() *Sensor {
	return &Sensor{
		Pause:    time.Second,
		Last:     time.Now(),
		modified: make([]string, 0),
		ignore:   []string{"#", ".git/", "vendor/"},
		React:    noReact,
	}
}

func noReact(*Directory, ...string) {}

type Sensor struct {
	Recursive bool
	Pause     time.Duration // between scans
	Last      time.Time
	modified  []string
	ignore    []string
	root      string
	React     ModifiedFunc
}

func (s *Sensor) scanForChanges(root string) {
	s.root = root // set it so the visitor knows to enter the first dir
	filepath.Walk(root, s.visit)
}

// Ignore returns true if the file should be ignored
func (w *Sensor) Ignore(path string, f os.FileInfo) bool {
	for _, str := range w.ignore {
		if strings.Contains(path, str) {
			return true
		}
	}
	return false
}

func (w *Sensor) visit(path string, f os.FileInfo, err error) error {
	if w.Ignore(path, f) {
		if f.IsDir() {
			if w.Recursive {
				return nil
			}
			return filepath.SkipDir
		}
		return nil
	}
	if !f.IsDir() && f.ModTime().After(w.Last) {
		w.modified = append(w.modified, path)
	}
	if !f.IsDir() {
		return nil
	}
	if w.root == path {
		// the starting directory
		return nil
	}
	if w.Recursive {
		return nil
	}
	return filepath.SkipDir
}

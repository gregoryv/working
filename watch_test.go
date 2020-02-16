package working

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wd, _ := TempDir()
	wd.MkdirAll("sub", "vendor/a/b")
	var (
		called bool
		sens   = NewSensor()
	)
	sens.React = func(*Directory, ...string) { called = true }
	sens.Pause = 50 * time.Millisecond
	plus := sens.Pause + 10*time.Millisecond
	go wd.Watch(ctx, sens)
	time.Sleep(plus)

	shouldSense := func(s string, err error) {
		t.Helper()
		called = false
		time.Sleep(plus)
		if !called {
			t.Error(s)
		}
	}
	shouldSense(wd.Touch("a"))

	shouldNotSense := func(s string, err error) {
		t.Helper()
		called = false
		time.Sleep(plus)
		if called {
			t.Error(s)
		}
	}
	// Not recursive
	shouldNotSense(wd.Touch("sub/hello"))

	// vendor should be ignored by default
	shouldNotSense(wd.Touch("vendor/noop"))

	// Directories are ignored by default
	shouldNotSense(wd.Touch("vendor"))

	// Removed
	sens.Recursive = true
	wd.MkdirAll("sub")
	wd.Touch("sub/x")
	os.RemoveAll(wd.Join("sub"))
	time.Sleep(plus)
	shouldNotSense("", nil)

	cancel()
	wd.RemoveAll()
}

package working

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatch_long(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	d := new(Directory)
	d.Temporary()
	defer d.RemoveAll()

	var (
		calls int
		sens  = NewSensor()
	)
	sens.Pause = 50 * time.Millisecond
	plus := sens.Pause + 10*time.Millisecond

	sens.React = func(*Directory, ...string) { calls++ }
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go d.Watch(ctx, sens)
	time.Sleep(plus)
	d.Touch("x")
	time.Sleep(plus)
	d.Touch("x")
	time.Sleep(5 * plus)
	if calls != 2 {
		t.Errorf("File changed twice but sensor reacted %v times", calls)
	}
}

func TestWatch(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	defer d.RemoveAll()
	d.MkdirAll("sub", "vendor/a/b")
	var (
		called bool
		sens   = NewSensor()
	)
	sens.React = func(*Directory, ...string) { called = true }
	sens.Pause = 50 * time.Millisecond
	plus := sens.Pause + 10*time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	go d.Watch(ctx, sens)
	defer cancel()
	time.Sleep(plus)

	shouldSense := func(s string, err error) {
		t.Helper()
		called = false
		time.Sleep(plus)
		if !called {
			t.Error(s)
		}
	}
	shouldSense(d.Touch("a"))

	shouldNotSense := func(s string, err error) {
		t.Helper()
		called = false
		time.Sleep(plus)
		if called {
			t.Error(s)
		}
	}
	// Not recursive
	shouldNotSense(d.Touch("sub/hello"))

	// vendor should be ignored by default
	shouldNotSense(d.Touch("vendor/noop"))

	// Directories are ignored by default
	shouldNotSense(d.Touch("vendor"))

	// Removed
	sens.Recursive = true
	d.MkdirAll("sub")
	d.Touch("sub/x")
	os.RemoveAll(d.Join("sub"))
	time.Sleep(plus)
	shouldNotSense("", nil)
}

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
		calls    int
		multiple bool
		sens     = NewSensor()
	)
	sens.Pause = 100 * time.Millisecond
	plus := 3 * sens.Pause

	sens.React = func(d *Directory, modified ...string) {
		calls++
		if len(modified) > 1 {
			multiple = true
		}
		d.Touch("y")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go d.Watch(ctx, sens)
	time.Sleep(100 * time.Millisecond)
	d.Touch("x")
	time.Sleep(plus)
	d.Touch("x")
	time.Sleep(plus)
	if calls != 2 {
		t.Errorf("File changed twice but sensor reacted %v times", calls)
	}
	if multiple {
		t.Error("Got multiple changes")
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

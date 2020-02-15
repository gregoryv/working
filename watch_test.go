package workdir

import (
	"context"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wd, _ := TempDir()
	var (
		called bool
		sens   = NewSensor()
	)
	go wd.Watch(ctx, sens, func(in WorkDir, modified ...string) {
		called = true
	})
	time.Sleep(40 * time.Millisecond)
	wd.Touch("hello")
	time.Sleep(1010 * time.Millisecond)
	if !called {
		t.Fail()
	}
	cancel()
	wd.RemoveAll()
}

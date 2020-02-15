/* package defines AssertFunc wrapper of testing.T

Online assertions are done by wrapping the T in a test

    func TestSomething(t *testing.T) {
        assert := asserter.New(t)
        got, err := Something()
        t.Logf("%v, %v := Something()", got, err)
        assert(err == nil).Fail()
        // Special case used very often is check equality
        assert().Equals(got, 1)
    }
*/
package asserter

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
)

type T interface {
	Helper()
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fail()
	FailNow()
	Log(...interface{})
	Logf(string, ...interface{})
}

type A interface {
	T
	Equals(got, exp interface{}) T
	Contains(body, exp interface{}) T
}

type AssertFunc func(expr ...bool) A

type wrappedT struct {
	T
}

func (w *wrappedT) Helper() {
	/* Cannot use the asserter as helper */
}
func (w *wrappedT) Error(args ...interface{}) {
	w.T.Helper()
	w.T.Error(args...)
}

func (w *wrappedT) Errorf(format string, args ...interface{}) {
	w.T.Helper()
	w.T.Errorf(format, args...)
}

func (w *wrappedT) Fatal(args ...interface{}) {
	w.T.Helper()
	w.T.Fatal(args...)
}

func (w *wrappedT) Fatalf(format string, args ...interface{}) {
	w.T.Helper()
	w.T.Fatalf(format, args...)
}

func (w *wrappedT) Fail() {
	w.T.Helper()
	w.T.Fail()
}

func (w *wrappedT) FailNow() {
	w.T.Helper()
	w.T.FailNow()
}
func (w *wrappedT) Log(args ...interface{}) {
	w.T.Helper()
	w.T.Log(args...)
}

func (w *wrappedT) Logf(format string, args ...interface{}) {
	w.T.Helper()
	w.T.Logf(format, args...)
}

// Helpers

func (w *wrappedT) Equals(got, exp interface{}) T {
	w.T.Helper()
	if got != exp {
		w.Errorf("got %v, expected %v", got, exp)
	}
	return w.T
}

func (w *wrappedT) Contains(body, exp interface{}) T {
	w.T.Helper()
	b := toBytes(w.T, body, "body")
	e := toBytes(w.T, exp, "exp")

	if bytes.Index(b, e) == -1 {
		format := "%q does not contain %q"
		if bytes.Index(b, []byte("\n")) > -1 {
			format = "%s\ndoes not contain\n%s"
		}
		w.Errorf(format, string(b), string(e))
	}
	return w.T
}

func toBytes(t T, v interface{}, name string) (b []byte) {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case int:
		return []byte(strconv.Itoa(v))
	case io.Reader:
		return bytesOrError(v)
	}
	t.Fatalf("%s must be io.Reader, []byte, string or int", name)
	return
}

func bytesOrError(r io.Reader) []byte {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte(err.Error())
	}
	return body
}

type noopT struct{}

func (t *noopT) Helper()                          {}
func (t *noopT) Error(...interface{})             {}
func (t *noopT) Errorf(string, ...interface{})    {}
func (t *noopT) Fatal(...interface{})             {}
func (t *noopT) Fatalf(string, ...interface{})    {}
func (t *noopT) Fail()                            {}
func (t *noopT) FailNow()                         {}
func (t *noopT) Log(...interface{})               {}
func (t *noopT) Logf(string, ...interface{})      {}
func (t *noopT) Equals(got, exp interface{}) T    { return t }
func (t *noopT) Contains(body, exp interface{}) T { return t }

var ok *noopT = &noopT{}

// Assert returns an asserter for online assertions.
func New(t T) AssertFunc {
	return func(expr ...bool) A {
		if len(expr) > 1 {
			t.Helper()
			t.Fatal("Only 0 or 1 bool expressions are allowed")
		}
		if len(expr) == 0 || !expr[0] {
			return &wrappedT{t}
		}
		return ok
	}
}

package dir

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	status := bytes.NewBufferString(``)
	_ = Parse(status.String())
}

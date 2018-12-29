package man

import (
	"os"
)

func init() {
	Usage("a", "something", true)
	Section("Options", "html")
	Usage("x", "word", 24)
	Set("header\n", "footer")
}

func ExampleWriteText() {
	WriteText(os.Stdout)
	// output:
	// header
	// a
	//   something (default true)
	//
	// Options
	// x
	//   word (default 24)
	// footer
}

func ExampleWriteHtml() {
	WriteHtml(os.Stdout)
	// output:
	// header
	// <dt><a name="a" href="#a">a</a></dt>
	// <dd>something (default true)</dd>
	// html
	// <dt><a name="x" href="#x">x</a></dt>
	// <dd>word (default 24)</dd>
	// footer
}

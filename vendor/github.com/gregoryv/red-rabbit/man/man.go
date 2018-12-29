// command line flag descriptions for rendering
package man

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

var (
	htmlTemplate *template.Template
	head, foot   string
)

func init() {
	var err error
	htmlTemplate, err = template.New("test").Parse(`<dt><a name="{{.Name}}" href="#{{.Name}}">{{.Flags}}</a></dt>
<dd>{{.Summary}} {{if .NotEmpty}}(default {{.DefaultValue}}){{end}}</dd>`)
	if err != nil {
		panic(err)
	}
}

type Element interface {
	WriteText(w io.Writer)
	WriteHtml(w io.Writer)
}

type section struct {
	text, html string
}

func (s *section) WriteText(w io.Writer) {
	fmt.Fprintf(w, "%s", s.text)
}

func (s *section) WriteHtml(w io.Writer) {
	fmt.Fprintf(w, "%s", s.html)
}

func Section(text, html string) {
	all = append(all, &section{text, html})
}

type help struct {
	Name         string // used for anchors in html output
	Flags        string
	DefaultValue interface{}
	Summary      string
}

func (h *help) NotEmpty() bool {
	switch v := h.DefaultValue.(type) {
	case string:
		return v != ""
	case bool:
		return v
	default:
		return true
	}
}

func (h *help) WriteText(w io.Writer) {
	fmt.Fprintf(w, "%s\n  %s (default %v)\n", h.Flags, h.Summary, h.DefaultValue)
}

func (h *help) WriteHtml(w io.Writer) {
	err := htmlTemplate.Execute(w, h)
	if err != nil {
		panic(err)
	}
}

func Value(flags string, summary string) {
	connect(flags, summary, false)
}

func Usage(flags string, summary string, defaultValue interface{}) {
	connect(flags, summary, defaultValue)
}

var all = make([]Element, 0)

func connect(flags string, summary string, defaultValue interface{}) {
	name := strings.Replace(flags, "-", "", -1)
	name = strings.Replace(name, ",", "_", -1)
	name = strings.Replace(name, " ", "", -1)
	h := &help{Name: name, Flags: flags, DefaultValue: defaultValue, Summary: summary}

	all = append(all, h)
}

func Set(header, footer string) {
	head = header
	foot = footer
}

func WriteText(w io.Writer) {
	if head != "" {
		fmt.Fprintf(w, "%s", head)
	}
	last := len(all) - 1
	for i, e := range all {
		e.WriteText(w)
		if i != last {
			fmt.Fprint(w, "\n")
		}
	}
	if foot != "" {
		fmt.Fprintf(w, "%s\n", foot)
	}
}

func WriteHtml(w io.Writer) {
	if head != "" {
		fmt.Fprintf(w, "%s", head)
	}
	for _, e := range all {
		e.WriteHtml(w)
		fmt.Fprint(w, "\n")
	}
	if foot != "" {
		fmt.Fprintf(w, "%s", foot)
	}
}

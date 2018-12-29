package text

import (
	"strings"
)

type Markdown struct {
	fontSize int32
}

func NewMarkdown() *Markdown {
	return &Markdown{fontSize: size}
}

func (m *Markdown) Render(txt string) {
	for _, line := range strings.Split(txt, "\n") {
		// Poor mans parsing of markdown, far from complete
		// expanded on a need to basis
		if line == "" {
			line = " "
		}
		switch true {
		case strings.Index(line, "# ") == 0:
			write(line[2:], "FreeSerif", gold(3, m.fontSize), 0)
		case strings.Index(line, "## ") == 0:
			write(line[3:], "FreeSerif", gold(2, m.fontSize), 0)
		case strings.Index(line, "### ") == 0:
			write(line[4:], "FreeSerif", gold(1, m.fontSize), 0)
		case strings.Index(line, "    ") == 0:
			write(line, "FreeMono", gold(1, m.fontSize), 0)
		default:
			write(line, "FreeSans", gold(0, size), 0)
		}
	}
}

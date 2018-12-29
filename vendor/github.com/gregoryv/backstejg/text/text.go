package text

import (
	"github.com/gregoryv/backstejg/act"
	"strings"
)

var (
	size, x, y, ident int32 = 18, size, gold(1, size), 0
	fontColor               = "999999"
	fontSize                = size
)

func SetPosition(xpos, ypos int32) {
	x = xpos
	y = ypos
}

func SetSize(s int32) {
	size = s
}

func SetFontColor(color string) {
	fontColor = color
}

func gold(min int, s int32) int32 {
	res := float32(s)
	if min <= 0 {
		return size
	}
	for {
		res = res * 1.61 // golden mean
		min--
		if min == 0 {
			break
		}
	}
	return int32(res)
}

func write(txt, font string, fs, ident int32) {
	a := &act.Event{
		Code:      act.NONE,
		Delay:     1,
		Text:      txt,
		FontColor: fontColor,
		FontSize:  int(fs),
		Font:      font,
		X:         x + ident,
		Y:         y,
	}
	y += fs + gold(0, fontSize)/2 // New line
	send(a)
}

func send(a *act.Event) {
	act.SendEvent(a, "localhost:9994")
}

type Plain struct {
	fontSize int32
	ident    int32
}

func NewPlain() *Plain {
	return &Plain{fontSize: size, ident: ident}
}

func (p *Plain) Render(txt string) {
	for _, line := range strings.Split(txt, "\n") {
		write(line, "FreeMono", p.fontSize, p.ident)
	}
}

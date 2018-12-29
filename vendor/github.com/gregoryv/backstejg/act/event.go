package act

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

const (
	ERR = int32(iota)
	NONE
	QUIT
	PAUSE
	CLEAR
	HIDE
	SAVE
	REDRAW
	REMOVE

	// Effects
	FADE_OUT
	PULSE
)

var NamedCommands = map[string]int32{
	"clear": CLEAR,
	"quit":  QUIT,
	"pause": PAUSE,
	"hide":  HIDE,
	"save":  SAVE,
	"rm":    REMOVE,
}

var NamedEffects = map[string]int32{
	"none":     NONE,
	"fade_out": FADE_OUT,
	"pulse":    PULSE,
}

var EventNames = map[int32]string{}

func init() {
	for name, code := range NamedCommands {
		EventNames[code] = name
	}
	for name, code := range NamedEffects {
		EventNames[code] = name
	}
}

type Event struct {
	Code      int32
	ImageURI  string
	Position  string
	X, Y      int32
	Delay     time.Duration
	FontColor string
	Font      string
	FontSize  int
	BgColor   string
	Text      string
	Multiply  int
	Points    []int  // Points are for geometric shapes, x1,y1,...,xN,yN
	DimSpeed  int    // 1..255
	Tag       string // Tagging events enables post manipulation
}

func NewEvent() *Event {
	return &Event{
		Code:      NONE,
		ImageURI:  "",
		Position:  "",
		X:         0,
		Y:         0,
		Delay:     time.Duration(0),
		FontColor: "white",
		FontSize:  34,
		BgColor:   "000000",
		Text:      "",
		DimSpeed:  3,
		Tag:       "",
	}
}

func (s *Event) String() string {
	return fmt.Sprintf("%s %s", EventNames[s.Code], s.ImageURI)
}

// SendEvent gob encodes the arguments to the bind address
func SendEvent(a *Event, bind string) error {
	conn, err := net.Dial("tcp", bind)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = gob.NewEncoder(conn).Encode(a)
	if err != nil {
		return err
	}
	return nil
}

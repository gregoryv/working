package main

import (
	"flag"
	"github.com/gregoryv/backstejg/act"
	"math"
	"time"
)

var fontColor = flag.String("fc", "999999", "font color")
var x = flag.Int("x", 0, "x position")
var y = flag.Int("y", 0, "y position")
var r = flag.Int("r", 30, "radius")

func main() {
	flag.Parse()
	for a := 450.0; a >= 90.0; a -= 20.0 {
		points := genPoints(*x, *y, float64(*r), a)
		drawLine(points)
		time.Sleep(30 * time.Millisecond)
	}
}

func genPoints(x, y int, radius, angle float64) []int {
	points := make([]int, 4)
	rad := angle * 3.14 / 180
	points[0] = x + int(math.Cos(rad)*radius)
	points[1] = y - int(math.Sin(rad)*radius)
	points[2] = x
	points[3] = y
	return points
}

func drawLine(points []int) {
	a := &act.Event{
		Code:      act.FADE_OUT,
		Delay:     0,
		DimSpeed:  28,
		FontColor: *fontColor,
		X:         int32(*x),
		Y:         int32(*y),
		Points:    points,
	}
	act.SendEvent(a, "localhost:9994")
}

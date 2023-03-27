package text

import (
	"fmt"

	"github.com/rockwell-uk/go-text/fonts"
)

type CharMetric struct {
	Char    string
	Width   float64
	Metrics fonts.GlyphMetrics
}

type LetterPosition struct {
	Char  string
	X     float64
	Y     float64
	Angle float64
}

type LineData struct {
	Pos    []float64
	Angle  float64
	Length float64
}

func (d LineData) String() string {
	return fmt.Sprintf("(%v, %v) %v:%v", d.Pos[0], d.Pos[1], d.Angle, d.Length)
}

type MultiLineData []LineData

func (ld MultiLineData) String() string {
	var r string

	for i, d := range ld {
		r += fmt.Sprintf("[%v] %+v\n", i, d)
	}

	return r
}

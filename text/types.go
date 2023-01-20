package text

import (
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

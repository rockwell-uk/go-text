package text

import (
	"strings"

	geos "github.com/twpayne/go-geos"

	"github.com/rockwell-uk/go-text/fonts"
)

var (
	gctx = geos.NewContext()
)

func SplitStringInTwo(s string, split func(string) bool) []string {
	if !split(s) {
		return []string{s}
	}

	n := strings.Count(s, " ")

	// if there are ony 2 words just split on the space
	if n == 1 {
		return strings.Split(s, " ")
	}

	l := len(s)

	// figure out positions of spaces
	sp := []int{}
	for i := 0; i < l; i++ {
		if string(s[i]) == " " {
			sp = append(sp, i)
		}
	}

	// decide which line each word belongs on
	dist := l / 2
	pos := 0
	for _, o := range sp {
		d := dist - o
		if d < 0 {
			d = -d
		}
		if d < dist {
			dist = d
			pos = o
		}
	}

	return []string{
		s[0:pos],
		s[pos+1:],
	}
}

func getFullWidth(gm fonts.GlyphMetrics) float64 {
	return gm.Advance
}

func getFullHeight(gm fonts.GlyphMetrics) float64 {
	return gm.Ascent
}

func ShouldSplit(s string) bool {
	// dont split short strings
	if len(s) < 12 {
		return false
	}

	n := strings.Count(s, " ")

	// if there arent any spaces dont split
	if n == 0 {
		return false
	}

	words := strings.Split(s, " ")

	// if there are 2 words
	// if the first or secord word is short dont split
	if len(words) == 2 {
		if len(words[0]) < 4 {
			return false
		}
		if len(words[1]) < 4 {
			return false
		}
	}

	return true
}

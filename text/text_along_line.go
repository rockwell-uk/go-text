package text

import (
	"fmt"
	"math"

	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/rockwell-uk/go-text/fonts"
)

type TextGlyph struct {
	Char     rune
	Pos      []float64
	Rotation float64
}

func TextAlongLine(gc *draw2dimg.GraphicContext, label string, lineCoords [][]float64, tf fonts.TypeFace) ([]TextGlyph, error) {
	charpositions, err := GetLetterPositions(label, lineCoords, tf)
	if err != nil {
		return []TextGlyph{}, err
	}

	textGlyphs := []TextGlyph{}

	for i, c := range label {
		x := charpositions[i].X
		y := charpositions[i].Y
		rotation := charpositions[i].Angle
		pos := []float64{
			x,
			y,
		}

		textGlyphs = append(textGlyphs, TextGlyph{Char: c, Pos: pos, Rotation: rotation})
	}

	return textGlyphs, nil
}

func GetLetterPositions(label string, lineCoords [][]float64, tf fonts.TypeFace) ([]LetterPosition, error) {
	charMetrics := getCharMetrics(label, tf)

	lineData := GetLineData(lineCoords)
	letterPositions := calculateLetterPositions(label, charMetrics, lineData, lineCoords, tf)

	numPositions := len(letterPositions)
	labelLength := len(label)

	if numPositions < labelLength {
		fm := fonts.GetFaceMetrics(tf)
		return letterPositions,
			fmt.Errorf("[%v] the letters dont fit on the line [%v:%v] (%v:%v)", label, numPositions, labelLength, fm.Height, tf.Spacing)
	}

	return letterPositions, nil
}

func calculateLetterPositions(label string, charMetrics []CharMetric, lineData []LineData, lineCoords [][]float64, tf fonts.TypeFace) []LetterPosition {
	var letterPositions []LetterPosition
	var charIndex int         // index of the current character
	var charsOnSegment int    // number of characters on the current segment
	var charMetric CharMetric // character metrics
	var remainder float64     // what is left of the current segment
	var nudge float64         // how far do we need to nudge the first char on the next segment
	var positionX, positionY float64
	var offsetX, offsetY float64

	fm := fonts.GetFaceMetrics(tf)

	charIndex = 0
	nudge = 0

	for s, line := range lineData {
		// reset segment count
		charsOnSegment = 0

		// starting point of the current line
		lineCoord := lineCoords[s]

		// offset required for each character to be centered on the line
		offsetX = math.Sin(line.Angle*(math.Pi/180)) * fm.Height / 3
		offsetY = math.Cos(line.Angle*(math.Pi/180)) * fm.Height / 3

		remainder = line.Length

		for i := charIndex; i < len(label); i++ {
			charMetric = charMetrics[charIndex]
			charWidth := charMetric.Width

			// calculate the x and y increase depending on the angle of the line
			xIncrease := math.Cos(line.Angle*(math.Pi/180)) * charWidth
			yIncrease := math.Sin(line.Angle*(math.Pi/180)) * charWidth

			// if this is the first loop in the current segment we need a starting point
			if charsOnSegment == 0 {
				positionX = lineCoord[0]
				positionY = lineCoord[1]
				if nudge != 0 {
					remainder -= nudge
					positionX += math.Cos(line.Angle*(math.Pi/180)) * nudge
					positionY += math.Sin(line.Angle*(math.Pi/180)) * nudge
					nudge = 0
				}
			}

			if remainder > charWidth/2 {
				letterPositions = append(letterPositions, LetterPosition{
					Char:  charMetric.Char,
					X:     positionX - offsetX,
					Y:     positionY + offsetY,
					Angle: line.Angle,
				})

				// move along the line
				positionX += xIncrease
				positionY += yIncrease

				// increase counts
				remainder -= charWidth
				charsOnSegment++
				charIndex++
			} else if charsOnSegment > 0 {
				nudge = -remainder
			}
		}
	}

	return letterPositions
}

func getCharMetrics(label string, tf fonts.TypeFace) []CharMetric {
	charMetrics := []CharMetric{}

	for i, r := range label {
		charMetrics = append(charMetrics, CharMetric{
			Char:    string(label[i]),
			Metrics: fonts.GetGlyphMetrics(tf, r),
			Width:   fonts.GetGlyphWidth(tf, r) + tf.Spacing,
		})
	}

	return charMetrics
}

// calculate the angle and distance travelled from the origin to each point along the line
// also record the coordintaes of each point from the origin.
func GetLineData(lineCoords [][]float64) MultiLineData {
	var lineData MultiLineData
	var prevPos []float64
	var x, y float64

	for key, pos := range lineCoords {
		if key > 0 {
			x = pos[0] - prevPos[0]
			y = pos[1] - prevPos[1]

			angle := math.Atan2(y, x) / (math.Pi / 180)
			hypotenuse := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))

			lineData = append(lineData, LineData{
				Pos: []float64{
					x,
					y,
				},
				Angle:  angle,
				Length: hypotenuse,
			})
		}

		prevPos = []float64{
			pos[0],
			pos[1],
		}
	}

	return lineData
}

package fonts

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"

	"github.com/rockwell-uk/go-text/fonts/ttf"
)

func TestGetTextWidth(t *testing.T) {

	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	pink := color.RGBA{0xEC, 0x74, 0xB4, 0xFF}

	tests := []struct {
		label    string
		fontSize float64
		expected []float64
	}{
		{
			"Mellor Street",
			float64(60),
			[]float64{58.48, 37.52, 25.11, 25.11, 41.21, 26.53, 33.34, 45.11, 21.62, 26.53, 37.52, 37.52, 21.62},
		},
		{
			"Pilsworth Road",
			float64(34),
			[]float64{26.69, 14.23, 14.23, 21.35, 26.62, 23.33, 15.04, 12.27, 25.47, 18.9, 27.21, 23.33, 21.27, 24.32},
		},
	}

	for _, tt := range tests {

		// Font
		f, err := truetype.Parse(ttf.ArialBold)
		if err != nil {
			t.Fatal(err)
		}

		// Truetype stuff
		opts := truetype.Options{
			Size: tt.fontSize,
		}
		face := truetype.NewFace(f, &opts)

		// strokestyle
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: 1.0,
		}

		backgroundStrokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: 1.0,
		}

		// font
		fontData := draw2d.FontData{
			Name:   "bold",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		}

		typeFace := TypeFace{
			StrokeStyle:           strokeStyle,
			Color:                 pink,
			Size:                  tt.fontSize,
			FontData:              fontData,
			Face:                  face,
			BackgroundColor:       pink,
			BackgroundStrokeStyle: backgroundStrokeStyle,
		}

		actual := []float64{}

		for i := 0; i < len(tt.label); i++ {

			char := string(tt.label[i])

			actual = append(actual, GetTextWidth(typeFace, char))
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%v: Expected [%+v], Got [%+v]", tt.label, tt.expected, actual)
		}
	}
}

package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path"
	"reflect"
	"testing"

	geos "github.com/rockwell-uk/go-geos"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/rockwell-uk/go-geom/geom"

	"github.com/rockwell-uk/go-text/fonts"
	"github.com/rockwell-uk/go-text/fonts/ttf"
)

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	pink  = color.RGBA{0xEC, 0x74, 0xB4, 0xFF}
)

func TestGetCharMetrics(t *testing.T) {

	tests := map[string]struct {
		wkt      string
		label    string
		fontSize float64
		spacing  float64
		expected []float64
	}{
		"Pilsworth Road": {
			"MULTILINESTRING((384342 409455.9999997657,384476.99999999994 409567.9999997657,384504.99999999994 409597.9999997654,384563.00000000006 409669.9999997661))",
			"Pilsworth Road",
			float64(34),
			float64(0),
			[]float64{20.77, 9.45, 9.45, 17, 30.23, 20.77, 13.23, 13.23, 20.77, 9.45, 22.67, 20.77, 18.91, 20.77},
		},
		"Mellor Street": {
			"MULTILINESTRING ((388874 413258.9999997683,388844.99999999994 413290.99999976775,388740.99999999994 413427.9999997701,388659.00000000006 413499.99999976833,388648 413512.9999997685,388583.00000000006 413820.9999997686))",
			"Mellor Street",
			float64(34),
			float64(0),
			[]float64{32.09, 18.91, 9.45, 9.45, 20.77, 13.23, 9.45, 22.67, 13.23, 13.23, 18.91, 18.91, 13.23},
		},
	}

	for name, tt := range tests {

		// Font
		f, err := truetype.Parse(ttf.UniversBold)
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

		typeFace := fonts.TypeFace{
			StrokeStyle:           strokeStyle,
			Color:                 pink,
			Size:                  tt.fontSize,
			FontData:              fontData,
			Face:                  face,
			BackgroundColor:       pink,
			BackgroundStrokeStyle: backgroundStrokeStyle,
			Spacing:               tt.spacing,
		}

		metrics := getCharMetrics(tt.label, typeFace)
		actual := []float64{}

		for _, m := range metrics {
			actual = append(actual, m.Width)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("@%v@ %v: Expected [%+v]\nActual [%+v]", name, tt.label, tt.expected, actual)
		}
	}
}

// nolint:dupl
func TestGetLetterPositions(t *testing.T) {

	tests := map[string]struct {
		wkt      string
		label    string
		fontSize float64
		spacing  float64
		zoom     float64
		expected []LetterPosition
	}{
		"Mellor Street": {
			"MULTILINESTRING((388874 413258.9999997683,388844.99999999994 413290.99999976775,388740.99999999994 413427.9999997701,388659.00000000006 413499.99999976833,388648 413512.9999997685,388583.00000000006 413820.9999997686))",
			"Mellor Street",
			float64(60),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "M", X: 388859.1802609714, Y: 413245.56961127336, Angle: 132.18444331631272},
				{Char: "e", X: 388820.9410296373, Y: 413289.6155843894, Angle: 127.20300888929643},
				{Char: "l", X: 388800.7581149327, Y: 413316.20269318344, Angle: 127.20300888929643},
				{Char: "l", X: 388790.67875039927, Y: 413329.4803176171, Angle: 127.20300888929643},
				{Char: "o", X: 388780.59938586585, Y: 413342.7579420508, Angle: 127.20300888929643},
				{Char: "r", X: 388758.43324886553, Y: 413371.9575648306, Angle: 127.20300888929643},
				{Char: " ", X: 388744.3209292368, Y: 413390.5478320341, Angle: 127.20300888929643},
				{Char: "S", X: 388739.2022372033, Y: 413402.96296639036, Angle: 138.7152891060774},
				{Char: "t", X: 388709.1295993965, Y: 413429.36820934206, Angle: 138.7152891060774},
				{Char: "r", X: 388691.5909845436, Y: 413444.7679687247, Angle: 138.7152891060774},
				{Char: "e", X: 388674.05236969073, Y: 413460.16772810736, Angle: 138.7152891060774},
				{Char: "e", X: 388646.4531998187, Y: 413483.865515681, Angle: 130.23635830904382},
				{Char: "t", X: 388625.9245866997, Y: 413520.7468649566, Angle: 101.91678096578067},
			},
		},
		"Pilsworth Road 1": {
			"MULTILINESTRING((384342 409455.9999997657,384476.99999999994 409567.9999997657,384504.99999999994 409597.9999997654,384563.00000000006 409669.9999997661))",
			"Pilsworth Road",
			float64(34),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "P", X: 384334.76365949906, Y: 409464.72237447667, Angle: 39.68010608220529},
				{Char: "i", X: 384350.7486938591, Y: 409477.9840326124, Angle: 39.68010608220529},
				{Char: "l", X: 384358.0216151255, Y: 409484.01786358893, Angle: 39.68010608220529},
				{Char: "s", X: 384365.29453639186, Y: 409490.05169456545, Angle: 39.68010608220529},
				{Char: "w", X: 384378.37809845834, Y: 409500.9062053169, Angle: 39.68010608220529},
				{Char: "o", X: 384401.6437502977, Y: 409520.20807943546, Angle: 39.68010608220529},
				{Char: "r", X: 384417.62878465775, Y: 409533.4697375712, Angle: 39.68010608220529},
				{Char: "t", X: 384427.81087443064, Y: 409541.9171009383, Angle: 39.68010608220529},
				{Char: "h", X: 384437.9929642035, Y: 409550.36446430546, Angle: 39.68010608220529},
				{Char: " ", X: 384453.97799856355, Y: 409563.6261224412, Angle: 39.68010608220529},
				{Char: "R", X: 384461.16762159235, Y: 409567.6467770428, Angle: 46.97493401060472},
				{Char: "o", X: 384476.6357763281, Y: 409584.2197999738, Angle: 46.97493401060472},
				{Char: "a", X: 384490.74300740194, Y: 409598.36766078236, Angle: 51.14662565986203},
				{Char: "d", X: 384502.60580896307, Y: 409613.09389720316, Angle: 51.14662565986203},
			},
		},
		"Pilsworth Road 2": {
			"POLYGON((300 200,298.078528040323 180.49096779838717,292.3879532511287 161.73165676349103,283.14696123025453 144.44297669803979,270.71067811865476 129.28932188134524,255.55702330196021 116.85303876974548,238.268343236509 107.61204674887132,219.50903220161283 101.92147195967695,200 100,180.49096779838717 101.92147195967695,161.73165676349103 107.61204674887132,144.4429766980398 116.85303876974545,129.28932188134524 129.28932188134524,116.85303876974547 144.44297669803979,107.61204674887132 161.73165676349106,101.92147195967695 180.49096779838723,100 200.00000000000009,101.92147195967698 219.5090322016129,107.61204674887136 238.26834323650908,116.85303876974555 255.55702330196033,129.28932188134536 270.71067811865487,144.44297669803993 283.14696123025465,161.7316567634912 292.38795325112875,180.4909677983874 298.0785280403231,200.00000000000026 300,219.50903220161308 298.078528040323,238.26834323650925 292.3879532511286,255.5570233019605 283.14696123025436,270.710678118655 270.7106781186545,283.1469612302547 255.55702330195993,292.3879532511288 238.26834323650863,298.07852804032314 219.50903220161243,300 200))",
			"Pilsworth Road",
			float64(34),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "P", X: 311.2787602356182, Y: 198.88913907626497, Angle: -95.62500000000007},
				{Char: "i", X: 308.5852138879073, Y: 176.08473505124144, Angle: -106.87499999999989},
				{Char: "l", X: 305.84202368785265, Y: 167.04164887857206, Angle: -106.87499999999989},
				{Char: "s", X: 302.1647364056044, Y: 155.98070398699963, Angle: -118.125},
				{Char: "w", X: 293.2655282896111, Y: 138.90764820220244, Angle: -129.375},
				{Char: "o", X: 271.3404795420657, Y: 115.14489279581252, Angle: -140.62500000000003},
				{Char: "r", X: 252.38645744752372, Y: 102.30760455744638, Angle: -151.87499999999994},
				{Char: "t", X: 238.42001586474572, Y: 95.814754199008, Angle: -163.125},
				{Char: "h", X: 223.69899114809584, Y: 90.94597640939344, Angle: -174.37500000000003},
				{Char: " ", X: 200.8072825276443, Y: 88.53231912401506, Angle: 174.37500000000003},
				{Char: "R", X: 191.40278686059204, Y: 89.45858110012941, Angle: 174.37500000000003},
				{Char: "o", X: 156.3891604127964, Y: 97.61693908625664, Angle: 151.87500000000003},
				{Char: "a", X: 136.35141384453564, Y: 108.8323190321136, Angle: 140.62500000000003},
				{Char: "d", X: 120.22837745082033, Y: 122.46527648353494, Angle: 129.375},
			},
		},
		"Turf Hill Road": {
			"MULTILINESTRING((390902 411492.9999997673,390951.00000000006 411523.99999976787,391010 411571.999999767,391052 411608.9999997665,391092.99999999994 411656.99999976583))",
			"Turf Hill Road",
			float64(34),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "T", X: 390895.94072725705, Y: 411502.57755990914, Angle: 32.319616508635505},
				{Char: "u", X: 390913.4930146818, Y: 411513.68206828006, Angle: 32.319616508635505},
				{Char: "r", X: 390931.04530210653, Y: 411524.786576651, Angle: 32.319616508635505},
				{Char: "f", X: 390941.35550297465, Y: 411530.7638687093, Angle: 39.13039955650276},
				{Char: " ", X: 390950.14431629435, Y: 411537.9140897151, Angle: 39.13039955650276},
				{Char: "H", X: 390957.4747916581, Y: 411543.8778662821, Angle: 39.13039955650276},
				{Char: "i", X: 390976.5185133703, Y: 411559.37106360705, Angle: 39.13039955650276},
				{Char: "l", X: 390983.84898873407, Y: 411565.33484017407, Angle: 39.13039955650276},
				{Char: "l", X: 390991.17946409783, Y: 411571.2986167411, Angle: 39.13039955650276},
				{Char: " ", X: 390998.5099394616, Y: 411577.2623933081, Angle: 39.13039955650276},
				{Char: "R", X: 391005.4032478821, Y: 411583.05436153454, Angle: 41.37851529552499},
				{Char: "o", X: 391022.4138862951, Y: 411598.0399239458, Angle: 41.37851529552499},
				{Char: "a", X: 391037.748043557, Y: 411609.76448261284, Angle: 49.49715161429619},
				{Char: "d", X: 391050.0298209829, Y: 411624.1431488674, Angle: 49.49715161429619},
			},
		},
	}

	for name, tt := range tests {

		geom, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		// Font
		f, err := truetype.Parse(ttf.UniversBold)
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

		typeFace := fonts.TypeFace{
			StrokeStyle:           strokeStyle,
			Color:                 pink,
			Size:                  tt.fontSize,
			FontData:              fontData,
			Face:                  face,
			BackgroundColor:       pink,
			BackgroundStrokeStyle: backgroundStrokeStyle,
			Spacing:               tt.spacing,
		}

		actual, g, err := GetLetterPositions(tt.label, geom, typeFace)
		if err != nil {
			t.Fatalf("%v %s", err, g)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%v: Expected [%+v]\nActual [%+v]", name, tt.expected, actual)
		}
	}
}

func TestLoveHeart(t *testing.T) {

	tests := map[string]struct {
		dim            int
		fontData       draw2d.FontData
		fontSize       float64
		fontStroke     float64
		fontSpacing    float64
		zoom           float64
		text           string
		heartSize      float64
		heartLineWidth float64
	}{
		"4000x4000": {
			dim: 4000,
			fontData: draw2d.FontData{
				Name:   "bold",
				Family: draw2d.FontFamilySans,
				Style:  draw2d.FontStyleBold,
			},
			fontSize:       60,
			zoom:           float64(1),
			fontStroke:     float64(10),
			fontSpacing:    float64(0.5),
			text:           "Why are there two equation expressions? Because the equation expression corresponding to the horizontal and vertical direction of the cardioid is different, and the cardioid drawn with the same equation expression will change the direction by exchanging the X coordinate and Y coordinate of each point, so there will be two equation expressions.",
			heartSize:      float64(100),
			heartLineWidth: 10,
		},
		"400x400": {
			dim: 400,
			fontData: draw2d.FontData{
				Name:   "bold",
				Family: draw2d.FontFamilySans,
				Style:  draw2d.FontStyleNormal,
			},
			fontSize:       6,
			zoom:           float64(1),
			fontStroke:     float64(1),
			fontSpacing:    float64(0.5),
			text:           "Why are there two equation expressions? Because the equation expression corresponding to the horizontal and vertical direction of the cardioid is different, and the cardioid drawn with the same equation expression will change the direction by exchanging the X coordinate and Y coordinate of each",
			heartSize:      float64(10),
			heartLineWidth: 1,
		},
	}

	for name, tt := range tests {

		tileWidth := float64(tt.dim)
		tileHeight := float64(tt.dim)
		centerX := tileWidth / 2
		centerY := tileHeight / 2

		bounds, err := geom.BoundsGeom(
			0,
			tileWidth,
			tileHeight,
			0,
		)
		if err != nil {
			t.Fatal(err)
		}

		envelope, err := geom.ToEnvelope(bounds)
		if err != nil {
			t.Fatal(err)
		}

		scale := func(x, y float64) (float64, float64) {
			nx := envelope.Px(x) * tileWidth
			ny := tileHeight - (envelope.Py(y) * tileHeight)
			return nx, ny
		}

		m := image.NewRGBA(image.Rect(0, 0, tt.dim, tt.dim))
		draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
		gc := draw2dimg.NewGraphicContext(m)

		gc.SetDPI(72)

		var wkt string
		var polyGeom *geos.Geom

		// heart shape
		wkt = generateHeart(tt.heartSize)
		polyGeom, err = gctx.NewGeomFromWKT(wkt)
		if err != nil {
			t.Fatal(err)
		}

		polyGeom = polyGeom.Simplify(2)

		// draw the line
		fillColour := white
		strokeColour := pink
		gc.Translate(centerX, centerY)
		err = geom.DrawLine(gc, polyGeom, tt.heartLineWidth, fillColour, tt.heartLineWidth, strokeColour, scale)
		if err != nil {
			t.Fatal(err)
		}

		fontSize := tt.fontSize * tt.zoom
		fontSpacing := tt.fontSpacing * tt.zoom
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: tt.fontStroke,
		}

		// font options
		face := fonts.GetFace(gc, tt.fontData, fontSize)

		typeFace := fonts.TypeFace{
			StrokeStyle: strokeStyle,
			Color:       pink,
			Size:        tt.fontSize,
			FontData:    tt.fontData,
			Face:        face,
			Spacing:     fontSpacing,
		}
		fonts.SetFont(gc, typeFace)

		// text along line
		gc.Translate(0, 0)
		glyphs, err := TextAlongLine(gc, tt.text, polyGeom, typeFace)
		if err != nil {
			t.Fatalf("%v: %v", name, err)
		}
		for _, glyph := range glyphs {
			err = geom.DrawRune(gc, glyph.Pos, face, glyph.Rotation, glyph.Char)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = savePNG(fmt.Sprintf("test-output/heart-test/%v.png", tt.dim), m)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTestAlongLineOutlines(t *testing.T) {

	tests := map[string]struct {
		dim         int
		fontData    draw2d.FontData
		fontSize    float64
		fontStroke  float64
		fontSpacing float64
		zoom        float64
		text        string
		lineWidth   float64
	}{
		"400x400": {
			dim: 400,
			fontData: draw2d.FontData{
				Name:   "bold",
				Family: draw2d.FontFamilySans,
				Style:  draw2d.FontStyleNormal,
			},
			fontSize:    40,
			zoom:        float64(1),
			fontStroke:  float64(1),
			fontSpacing: float64(0.5),
			text:        "Why are there two equation expressions?",
			lineWidth:   1,
		},
	}

	for _, tt := range tests {

		tileWidth := float64(tt.dim)
		tileHeight := float64(tt.dim)
		centerX := tileWidth / 2
		centerY := tileHeight / 2

		bounds, err := geom.BoundsGeom(
			0,
			tileWidth,
			tileHeight,
			0,
		)
		if err != nil {
			t.Fatal(err)
		}

		envelope, err := geom.ToEnvelope(bounds)
		if err != nil {
			t.Fatal(err)
		}

		scale := func(x, y float64) (float64, float64) {
			nx := envelope.Px(x) * tileWidth
			ny := tileHeight - (envelope.Py(y) * tileHeight)
			return nx, ny
		}

		m := image.NewRGBA(image.Rect(0, 0, tt.dim, tt.dim))
		draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
		gc := draw2dimg.NewGraphicContext(m)

		gc.SetDPI(72)
		gc.Translate(-200, -200)

		// generate a circular line to use for the test
		var polyGeom *geos.Geom

		radius := 140.0
		numPoints := 40 // IMPORTANT: the text will never fit on the line unless the distance between each point is greater than the charWidth/2
		origin := []float64{
			200.00,
			200.00,
		}
		polyGeom, err = geom.CircleGeom(
			origin,
			radius,
			numPoints,
		)
		if err != nil {
			t.Fatal(err)
		}

		// draw the line
		fillColour := black
		strokeColour := white
		gc.Translate(centerX, centerY)
		err = geom.DrawLine(gc, polyGeom, tt.lineWidth, fillColour, tt.lineWidth, strokeColour, scale)
		if err != nil {
			t.Fatal(err)
		}

		fontSize := tt.fontSize * tt.zoom
		fontSpacing := tt.fontSpacing * tt.zoom
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: tt.fontStroke,
		}

		// font options
		face := fonts.GetFace(gc, tt.fontData, fontSize)

		typeFace := fonts.TypeFace{
			StrokeStyle: strokeStyle,
			Color:       pink,
			Size:        tt.fontSize,
			FontData:    tt.fontData,
			Face:        face,
			Spacing:     fontSpacing,
		}
		fonts.SetFont(gc, typeFace)

		// text along line
		glyphs, err := TextAlongLine(gc, tt.text, polyGeom, typeFace)
		if err != nil {
			t.Fatal(err)
		}
		for _, glyph := range glyphs {
			err = geom.DrawRune(gc, glyph.Pos, face, glyph.Rotation, glyph.Char)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = DrawGlyphOutlines(gc, tt.text, polyGeom, typeFace)
		if err != nil {
			t.Fatal(err)
		}

		err = savePNG(fmt.Sprintf("test-output/outline-test/%v.png", tt.dim), m)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// https://developpaper.com/a-romantic-and-sad-love-story-cartesian-heart-line/
func generateHeart(size float64) string {

	// an array that holds the coordinates of all points
	var p [][]float64

	// t for radian
	t := float64(0)

	// vt represents the increment of T
	vt := 0.01

	// maxt represents the maximum value of T
	maxt := 2 * math.Pi

	// number of cycles required
	maxi := int(math.Ceil(maxt / vt))

	// x is used to temporarily save the X coordinate of each cycle
	var x float64

	// y is used to temporarily save the Y coordinate of each cycle
	var y float64

	// get the coordinates of all points according to the equation
	for i := 0; i <= maxi; i++ {
		x = 16 * math.Pow(math.Sin(t), 3)
		y = 13*math.Cos(t) - 5*math.Cos(2*t) - 2*math.Cos(3*t) - math.Cos(4*t)
		t += vt
		p = append(p, []float64{x * size, -y * size})
	}

	return toLineString(p)
}

func toLineString(p [][]float64) string {

	s := "LINESTRING("

	for i, pt := range p {
		s += fmt.Sprintf("%v %v", pt[0], pt[1])
		if i+1 < len(p) {
			s += ","
		}
	}
	s += ")"

	return s
}

func savePNG(fname string, m image.Image) error {

	dir, _ := path.Split(fname)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	return draw2dimg.SaveToPngFile(fname, m)
}

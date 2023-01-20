package fonts

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/rockwell-uk/csync/mutex"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type TypeFace struct {
	Name                  string
	Size                  float64
	Color                 color.RGBA
	BackgroundColor       color.RGBA
	BackgroundStrokeStyle draw2d.StrokeStyle
	FontData              draw2d.FontData
	Spacing               float64
	Face                  font.Face
	StrokeStyle           draw2d.StrokeStyle
}

type GlyphMetrics struct {
	Ascent       float64
	Descent      float64
	BearingLeft  float64
	BearingRight float64
	Advance      float64
}

type FaceMetrics struct {
	Ascent     float64
	Descent    float64
	Height     float64
	XHeight    float64
	CapHeight  float64
	CaretSlope image.Point
}

type GlyphBounds struct {
	BlX float64
	BlY float64
	TlX float64
	TlY float64
	TrX float64
	TrY float64
	BrX float64
	BrY float64
}

var (
	cache_glyphmetrics = make(map[string]GlyphMetrics)
	cache_glyphbounds  = make(map[string]GlyphBounds)
	cache_facemetrics  = make(map[string]FaceMetrics)
)

func GetTextWidth(tf TypeFace, text string) float64 {

	var w float64
	for i, char := range text {
		w += GetGlyphWidth(tf, char)

		if i+1 == len(text) {
			b := GetGlyphMetrics(tf, char)
			w += b.BearingRight + b.BearingLeft
		}
	}

	return math.Round(w*100) / 100
}

func GetGlyphWidth(tf TypeFace, char rune) float64 {

	b := GetGlyphMetrics(tf, char)
	adv := b.Advance

	return math.Round(adv*100) / 100
}

func GetTextHeight(tf TypeFace, text string) (float64, float64) {

	var maxAscent, maxDescent float64

	for _, char := range text {
		b := GetGlyphMetrics(tf, char)
		if b.Ascent > maxAscent {
			maxAscent = b.Ascent
		}
		if b.Descent > maxDescent {
			maxDescent = b.Descent
		}
	}

	return math.Round((maxAscent)*100) / 100, math.Round((maxDescent)*100) / 100
}

// GlyphBounds returns the bounding box of r's glyph, drawn at a dot equal
// to the origin, and that glyph's advance width.
//
// It returns !ok if the face does not contain a glyph for r.
//
// The glyph's ascent and descent are equal to -bounds.Min.Y and
// +bounds.Max.Y. The glyph's left-side and right-side bearings are equal
// to bounds.Min.X and advance-bounds.Max.X. A visual depiction of what
// these metrics are is at
// https://developer.apple.com/library/archive/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyphterms_2x.png
func GetGlyphMetrics(tf TypeFace, char rune) GlyphMetrics {

	cacheKey, ok := glyphCacheKey(tf, char)
	if ok {
		mutex.Lock()
		cachedVersion, exists := cache_glyphmetrics[cacheKey]
		mutex.Unlock()
		if exists {
			return cachedVersion
		}
	}

	bounds, advance, _ := tf.Face.GlyphBounds(char)
	gm := GlyphMetrics{
		Ascent:       unfix(-bounds.Min.Y),
		Descent:      unfix(bounds.Max.Y),
		BearingLeft:  unfix(bounds.Min.X),
		BearingRight: unfix(advance - bounds.Max.X),
		Advance:      unfix(advance),
	}

	mutex.Lock()
	cache_glyphmetrics[cacheKey] = gm
	mutex.Unlock()

	return gm
}

func GetGlyphBounds(tf TypeFace, char rune) GlyphBounds {

	cacheKey, ok := glyphCacheKey(tf, char)
	if ok {
		mutex.Lock()
		cachedVersion, exists := cache_glyphbounds[cacheKey]
		mutex.Unlock()
		if exists {
			return cachedVersion
		}
	}

	bounds, _, _ := tf.Face.GlyphBounds(char)
	gb := GlyphBounds{
		BlX: unfix(bounds.Min.X),
		BlY: unfix(-bounds.Min.Y),
		TlX: unfix(bounds.Max.X),
		TlY: unfix(-bounds.Min.Y),
		TrX: unfix(bounds.Max.X),
		TrY: unfix(bounds.Max.Y),
		BrX: unfix(bounds.Max.X),
		BrY: unfix(-bounds.Min.Y),
	}

	mutex.Lock()
	cache_glyphbounds[cacheKey] = gb
	mutex.Unlock()

	return gb
}

// Metrics holds the metrics for a Face. A visual depiction is at
// https://developer.apple.com/library/mac/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png
func GetFaceMetrics(tf TypeFace) FaceMetrics {

	cacheKey := fmt.Sprintf("%s.%v", tf.Name, tf.Size)
	if cacheKey != "" {
		mutex.Lock()
		cachedVersion, exists := cache_facemetrics[cacheKey]
		mutex.Unlock()
		if exists {
			return cachedVersion
		}
	}

	m := tf.Face.Metrics()
	fm := FaceMetrics{
		Ascent:     unfix(m.Ascent),
		Descent:    unfix(m.Descent),
		Height:     unfix(m.Height),
		XHeight:    unfix(m.XHeight),
		CapHeight:  unfix(m.CapHeight),
		CaretSlope: m.CaretSlope,
	}

	mutex.Lock()
	cache_facemetrics[cacheKey] = fm
	mutex.Unlock()

	return fm
}

//https://github.com/fogleman/gg/blob/master/util.go
func unfix(x fixed.Int26_6) float64 {

	const shift, mask = 6, 1<<6 - 1

	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}

	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}

	return 0
}

// nolint:ireturn
func GetFace(gc *draw2dimg.GraphicContext, fontData draw2d.FontData, size float64) font.Face {

	font, err := gc.FontCache.Load(fontData)
	if err != nil {
		panic(err)
	}

	// Truetype stuff
	opts := truetype.Options{
		Size: size,
	}

	return truetype.NewFace(font, &opts)
}

func SetFont(gc *draw2dimg.GraphicContext, typeFace TypeFace) {

	font, err := gc.FontCache.Load(typeFace.FontData)
	if err != nil {
		panic(err)
	}

	gc.SetFont(font)
	gc.SetFontData(typeFace.FontData)
	gc.SetFontSize(typeFace.Size)
	gc.SetFillColor(typeFace.Color)

	if typeFace.StrokeStyle.Color != nil {
		gc.SetStrokeColor(typeFace.StrokeStyle.Color)
	}

	gc.SetLineWidth(typeFace.StrokeStyle.Width)
}

func glyphCacheKey(typeFace TypeFace, r rune) (string, bool) {

	if typeFace.Name == "" {
		return "", false
	}

	s := string(r)
	if s == "" {
		return "", false

	}

	return fmt.Sprintf("%s.%v.%s", typeFace.Name, typeFace.Size, string(r)), true
}

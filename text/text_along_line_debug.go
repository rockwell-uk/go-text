package text

import (
	"fmt"
	"image/color"
	"math"

	geos "github.com/twpayne/go-geos"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/rockwell-uk/go-geom/geom"

	"github.com/rockwell-uk/go-text/fonts"
)

func DrawGlyphOutlines(gc *draw2dimg.GraphicContext, label string, g *geos.Geom, tf fonts.TypeFace) error {

	var (
		black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
		white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	)

	letterpositions, _, _ := GetLetterPositions(label, g, tf)

	var blx, bly, tlx, tly, trx, try, brx, bry float64
	var tblx, tbly, ttlx, ttly, ttrx, ttry, tbrx, tbry float64
	var wkt string
	var polyGeom *geos.Geom

	scale := func(x, y float64) (float64, float64) {
		return x, y
	}

	if len(letterpositions) >= len(label) {

		for i, r := range label {

			x := letterpositions[i].X
			y := letterpositions[i].Y

			gm := fonts.GetGlyphMetrics(tf, r)

			rotation := letterpositions[i].Angle
			radians := rotation * (math.Pi / 180)

			blx = x
			bly = y + tf.StrokeStyle.Width
			tlx = blx
			tly = bly - getFullHeight(gm) - (tf.StrokeStyle.Width * 2)
			trx = tlx + getFullWidth(gm) + tf.StrokeStyle.Width
			try = tly
			brx = trx
			bry = try + getFullHeight(gm) + (tf.StrokeStyle.Width * 2)

			tblx, tbly = rotateAroundPoint(blx, bly, x, y, radians)
			ttlx, ttly = rotateAroundPoint(tlx, tly, x, y, radians)
			ttrx, ttry = rotateAroundPoint(trx, try, x, y, radians)
			tbrx, tbry = rotateAroundPoint(brx, bry, x, y, radians)

			wkt = fmt.Sprintf("LINESTRING(%v %v, %v %v, %v %v, %v %v, %v %v)", tblx, tbly, ttlx, ttly, ttrx, ttry, tbrx, tbry, tblx, tbly)

			polyGeom, _ = gctx.NewGeomFromWKT(wkt)

			err := geom.DrawLine(gc, polyGeom, 1.0, black, 1.0, white, scale)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func rotateAroundPoint(x, y, originx, originy, radians float64) (float64, float64) {

	rx := originx + (x-originx)*math.Cos(radians) - (y-originy)*math.Sin(radians)
	ry := originy + (x-originx)*math.Sin(radians) + (y-originy)*math.Cos(radians)

	return rx, ry
}

package text

import (
	"fmt"

	geos "github.com/twpayne/go-geos"
)

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

func (ld MultiLineData) ToWKT(origin []float64) string {

	var r string

	n := len(ld)

	if n == 0 {
		return "LINESTRING EMPTY"
	} else {
		r = "LINESTRING ("

		x := origin[0]
		y := origin[1]

		r = fmt.Sprintf("%v%v %v,", r, x, y)

		for i, p := range ld {

			x += p.Pos[0]
			y += p.Pos[1]
			r = fmt.Sprintf("%v%v %v", r, x, y)
			if i < n-1 {
				r = fmt.Sprintf("%v,", r)
			}
		}

		r = fmt.Sprintf("%v)", r)
	}

	return r
}

func (ld MultiLineData) ToGeom(origin []float64) (*geos.Geom, error) {

	r := ld.ToWKT(origin)

	g, err := gctx.NewGeomFromWKT(r)
	if err != nil {
		return &geos.Geom{}, err
	}

	return g, nil
}

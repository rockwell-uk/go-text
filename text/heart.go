package text

import "math"

// https://developpaper.com/a-romantic-and-sad-love-story-cartesian-heart-line/
func GenerateHeart(size float64) [][]float64 {
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

	return p
}

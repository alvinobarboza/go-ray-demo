package raytracer

import (
	"math"
)

const (
	TAU        = 2 * math.Pi
	DEG_TO_RAD = TAU / 360
)

func VecDot(v1, v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func VecLen(v Vec3) float64 {
	return math.Sqrt(VecDot(v, v))
}

func VecNormal(v Vec3) Vec3 {
	n := VecLen(v)
	return Vec3{
		X: v.X / n,
		Y: v.Y / n,
		Z: v.Z / n,
	}
}

func VecAdd(v1, v2 Vec3) Vec3 {
	return Vec3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func VecMultiply(v Vec3, n float64) Vec3 {
	return Vec3{
		X: v.X * n,
		Y: v.Y * n,
		Z: v.Z * n,
	}
}

func MatrixMultiplication(m []float64, v Vec3) Vec3 {
	result := []float64{0, 0, 0}
	vec := []float64{v.X, v.Y, v.Z}

	length := 3

	for h := range length {
		for w := range length {
			result[h] += vec[w] * m[length*h+w]
		}
	}
	return Vec3{X: result[0], Y: result[1], Z: result[2]}
}

func RotateXYZ(angle Vec3, v Vec3) Vec3 {
	cosa := math.Cos(angle.X * -DEG_TO_RAD)
	sina := math.Sin(angle.X * -DEG_TO_RAD)

	cosb := math.Cos(angle.Y * -DEG_TO_RAD)
	sinb := math.Sin(angle.Y * -DEG_TO_RAD)

	cosga := math.Cos(angle.Z * -DEG_TO_RAD)
	singa := math.Sin(angle.Z * -DEG_TO_RAD)

	// Formula for general 3D roation using matrix
	matrix := []float64{
		cosb * cosga, sina*sinb*cosga - cosa*singa, cosa*sinb*cosga + sina*singa,
		cosb * singa, sina*sinb*singa + cosa*cosga, cosa*sinb*singa - sina*cosga,
		-sinb, sina * cosb, cosa * cosb,
	}

	value := MatrixMultiplication(matrix, v)

	return value
}

func CrossProdutc(v, w Vec3) Vec3 {
	return Vec3{
		X: v.Y*w.Z - v.Z*w.Y,
		Y: v.Z*w.X - v.X*w.Z,
		Z: v.X*w.Y - v.Y*w.X,
	}
}
